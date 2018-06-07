package fmttab

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/arteev/fmttab/columns"
	"github.com/arteev/fmttab/eol"
)

func (t *Table) writeHeader(buf *bufio.Writer) (int, error) {
	if t.caption != "" {
		buf.WriteString(t.caption)
		buf.WriteString(eol.EOL)
	}
	buf.WriteString(Borders[t.border][BKLeftTop])

	cntCols := t.columnsvisible.Len()
	num := 0
	t.columnsvisible.Visit(func(c *columns.Column) error {
		cnw, _ := buf.WriteString(strings.Repeat(Borders[t.border][BKHorizontalBorder], c.GetWidth()))
		if num < cntCols-1 {
			buf.WriteString(Borders[t.border][BKTopToBottom])
		} else {
			buf.WriteString(Borders[t.border][BKRighttop])
			if cnw > 0 {
				buf.WriteString(eol.EOL)
			}
		}
		num++
		return nil
	})

	if t.VisibleHeader {
		buf.WriteString(Borders[t.border][BKVerticalBorder])
		num := 0
		t.columnsvisible.Visit(func(c *columns.Column) error {
			caption := fmt.Sprintf(c.GetMaskFormat(), c.Caption)
			buf.WriteString(trimEnds(caption, c.GetWidth()))
			bKind := BKVertical
			if num == cntCols-1 {
				bKind = BKVerticalBorder
			}
			buf.WriteString(Borders[t.border][bKind])
			num++
			return nil
		})
		buf.WriteString(eol.EOL)
		buf.WriteString(Borders[t.border][BKLeftToRight])
		s := t.getBorderTopButtomData(BKHorizontal, BKBottomCross, BKRightToLeft)
		if _, err := buf.WriteString(s); err != nil {
			return 0, err
		}
	}
	return buf.Buffered(), buf.Flush()
}

func (t *Table) getBorderTopButtomData(hr, vbwnCol, vright BorderKind) string {
	var result string
	count := t.columnsvisible.Len()
	num := 0
	t.columnsvisible.Visit(func(c *columns.Column) error {
		strLine := strings.Repeat(Borders[t.border][hr], c.GetWidth())
		if num < count-1 {
			strLine += Borders[t.border][vbwnCol]
		} else {
			strLine += Borders[t.border][vright]
			if strLine != "" {
				strLine += eol.EOL
			}
		}
		result += strLine
		num++
		return nil
	})
	return result
}

func (t *Table) writeBottomBorder(buf *bufio.Writer) (int, error) {
	s := Borders[t.border][BKLeftBottom] + t.getBorderTopButtomData(BKHorizontalBorder, BKBottomToTop, BKRightBottom)
	if _, err := buf.WriteString(s); err != nil {
		return 0, err
	}
	return buf.Buffered(), buf.Flush()
}

func (t *Table) writeRecord(data map[string]interface{}, buf *bufio.Writer) (int, error) {
	var cntwrite int

	cntCols := t.columnsvisible.Len()
	if n, err := buf.WriteString(Borders[t.border][BKVerticalBorder]); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}

	num := 0
	err := t.columnsvisible.Visit(func(c *columns.Column) error {
		val, mok := data[c.Name]
		if !mok || val == nil {
			val = ""
		}

		mask, ok := t.masks[c.Name]
		if !ok {
			mask = c.GetMaskFormat()
			t.masks[c.Name] = mask
		}

		caption := fmt.Sprintf(mask, val)
		n, err := buf.WriteString(trimEnds(caption, c.GetWidth()))
		if err != nil {
			return err
		}
		cntwrite += n

		if num < cntCols-1 {
			n, err = buf.WriteString(Borders[t.border][BKVertical])
		} else {
			n, err = buf.WriteString(Borders[t.border][BKVerticalBorder])
		}
		if err != nil {
			return err
		}
		cntwrite += n

		num++
		return nil
	})

	if err != nil {
		return -1, err
	}
	n, err := buf.WriteString(eol.EOL)
	if err != nil {
		return -1, err
	}
	cntwrite += n

	return cntwrite, nil
}

func (t *Table) getRecordHorBorder() string {
	var result string
	count := t.columnsvisible.Len()
	result += Borders[t.border][BKLeftToRight]
	for num, c := range t.columnsvisible.Columns() {
		strLine := strings.Repeat(Borders[t.border][BKHorizontal], c.GetWidth())
		if num < count-1 {
			strLine += Borders[t.border][BKBottomCross]
		} else {
			strLine += Borders[t.border][BKRightToLeft]
		}
		result += strLine
	}
	return result + eol.EOL
}

func (t *Table) writeData(buf *bufio.Writer) (int, error) {
	firstrow := true
	var recordSeparator string
	if t.dataget != nil {
		for {
			ok, data := t.dataget()
			if !ok {
				break
			}
			if (!firstrow) && t.CloseEachColumn {
				if recordSeparator == "" {
					recordSeparator = t.getRecordHorBorder()
				}
				if _, err := buf.WriteString(recordSeparator); err != nil {
					return -1, err
				}
			}
			firstrow = false
			if _, err := t.writeRecord(data, buf); err != nil {
				return -1, err
			}
		}
		return buf.Buffered(), buf.Flush()
	}

	if t.CountData() == 0 {
		return buf.Buffered(), buf.Flush()
	}
	for ii, data := range t.Data {
		if _, err := t.writeRecord(data, buf); err != nil {
			return -1, err
		}
		if t.CloseEachColumn {
			if ii < len(t.Data)-1 {
				if recordSeparator == "" {
					recordSeparator = t.getRecordHorBorder()
				}
				if _, err := buf.WriteString(recordSeparator); err != nil {
					return -1, err
				}
			}
		}
	}

	return buf.Buffered(), buf.Flush()
}

func (t *Table) adjustmentWidth() error {
	resized := false
	t.Columns.Visit(func(c *columns.Column) error {
		if t.autoSize > 0 || c.IsAutoSize() {
			c.MaxLen = len(c.Caption)
			resized = true

			//loop on data
			for _, data := range t.Data {
				curval := fmt.Sprintf("%v", data[c.Name])
				curlen := utf8.RuneCountInString(curval)
				if curlen > c.MaxLen {
					c.MaxLen = curlen
				}
			}

		}

		return nil
	})
	if !resized {
		return nil
	}
	//adjustment of table
	if t.autoSize > 0 {
		termwidth := t.autoSize - utf8.RuneCountInString(Borders[t.border][BKVertical])*t.columnsvisible.Len() - utf8.RuneCountInString(Borders[t.border][BKVerticalBorder])*2
		nowwidths := make(map[string]int, t.columnsvisible.Len())
		allcolswidth := 0

		t.columnsvisible.Visit(func(c *columns.Column) error {
			currentWidth := c.Width
			if c.MaxLen > c.Width || c.IsAutoSize() {
				currentWidth = c.MaxLen
			}
			nowwidths[c.Name] = currentWidth
			allcolswidth += currentWidth
			return nil
		})

		twAll := 0
		t.columnsvisible.Visit(func(c *columns.Column) error {
			c.MaxLen = int(math.Trunc(float64(termwidth) * (float64(nowwidths[c.Name]) / float64(allcolswidth))))
			twAll += c.MaxLen
			return nil
		})

		i := 0
		//distrib mod
		for {
			if twAll >= termwidth || twAll <= 0 {
				break
			}
			if i+1 >= t.columnsvisible.Len() {
				i = 0
			}
			t.columnsvisible.Get(i).MaxLen = t.columnsvisible.Get(i).MaxLen + 1

			twAll = twAll + 1
			i = i + 1
		}
	}
	return nil
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	t.masks = make(map[string]string)
	t.columnsvisible = t.Columns.ColumnsVisible()
	buf := bufio.NewWriter(w)
	if t.columnsvisible.Len() == 0 {
		return 0, nil
	}
	if err := t.adjustmentWidth(); err != nil {
		return 0, err
	}
	var cntwrite int64
	if n, err := t.writeHeader(buf); err == nil {
		cntwrite += int64(n)
	} else {
		return -1, err
	}
	if n, err := t.writeData(buf); err == nil {
		cntwrite += int64(n)
	} else {
		return -1, err
	}
	if n, err := t.writeBottomBorder(buf); err == nil {
		cntwrite += int64(n)
	} else {
		return -1, err
	}
	return cntwrite, nil
}

// String returns the contents of the table with borders
// as a string.  If error, it returns "".
func (t *Table) String() string {
	var buf bytes.Buffer
	if _, err := t.WriteTo(&buf); err != nil {
		return ""
	}
	return buf.String()
}
