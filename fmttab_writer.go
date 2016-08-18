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
	"github.com/arteev/tern"
)

func (t *Table) writeHeader(buf *bufio.Writer) (int, error) {
	if t.caption != "" {
		buf.WriteString(t.caption)
		buf.WriteString(eol.EOL)
	}
	buf.WriteString(Borders[t.border][BKLeftTop])
	colv := t.Columns.ColumnsVisible()
	cntCols := colv.Len()
	for num, c := range colv {
		cnw, _ := buf.WriteString(strings.Repeat(Borders[t.border][BKHorizontalBorder], t.getWidth(c)))
		if num < cntCols-1 {
			buf.WriteString(Borders[t.border][BKTopToBottom])
		} else {
			buf.WriteString(Borders[t.border][BKRighttop])
			if cnw > 0 {
				buf.WriteString(eol.EOL)
			}
		}
	}
	if t.VisibleHeader {
		buf.WriteString(Borders[t.border][BKVerticalBorder])
		for num, c := range colv {
			caption := fmt.Sprintf(t.GetMaskFormat(c), c.Caption)
			buf.WriteString(trimEnds(caption, Trimend, t.getWidth(c)))
			if num < cntCols-1 {
				buf.WriteString(Borders[t.border][BKVertical])
			} else {
				buf.WriteString(Borders[t.border][BKVerticalBorder])
			}
		}
		buf.WriteString(eol.EOL)
		buf.WriteString(Borders[t.border][BKLeftToRight])
		if err := t.writeBorderTopButtomData(buf, BKHorizontal, BKBottomCross, BKRightToLeft); err != nil {
			return 0, err
		}
	}
	return buf.Buffered(), buf.Flush()
}

func (t *Table) writeBorderTopButtomData(b *bufio.Writer, hr, vbwnCol, vright BorderKind) error {
	colv := t.Columns.ColumnsVisible()
	empty := true
	for num, c := range colv {
		cnt, err := b.WriteString(strings.Repeat(Borders[t.border][hr], t.getWidth(c)))
		if err != nil {
			return err
		}
		if cnt > 0 {
			empty = false
		}
		if num < colv.Len()-1 {
			cnt, err = b.WriteString(Borders[t.border][vbwnCol])
			if err != nil {
				return err
			}
		} else {
			cnt, err = b.WriteString(Borders[t.border][vright])
			if err != nil {
				return err
			}
			if !empty {
				ceol, _ := b.WriteString(eol.EOL)
				cnt += ceol
			}
		}
		if cnt > 0 {
			empty = false
		}
	}
	return nil
}

func (t *Table) writeBottomBorder(buf *bufio.Writer) (int, error) {
	if _, err := buf.WriteString(Borders[t.border][BKLeftBottom]); err != nil {
		return 0, err
	}
	if err := t.writeBorderTopButtomData(buf, BKHorizontalBorder, BKBottomToTop, BKRightBottom); err != nil {
		return 0, err
	}
	return buf.Buffered(), buf.Flush()
}

func (t *Table) writeRecord(data map[string]interface{}, buf *bufio.Writer) (int, error) {
	var cntwrite int
	colv := t.Columns.ColumnsVisible()
	cntCols := tern.Op(colv == nil, 0, colv.Len()).(int)
	if n, err := buf.WriteString(Borders[t.border][BKVerticalBorder]); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}
	for num, c := range colv {
		val, mok := data[c.Name]
		if !mok || val == nil {
			val = ""
		}
		caption := fmt.Sprintf(t.GetMaskFormat(c), val)
		if n, err := buf.WriteString(trimEnds(caption, Trimend, t.getWidth(c))); err == nil {
			cntwrite += n
		} else {
			return -1, err
		}
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKVertical]
		} else {
			delim = Borders[t.border][BKVerticalBorder]
		}
		if n, err := buf.WriteString(delim); err == nil {
			cntwrite += n
		} else {
			return -1, err
		}
	}
	if n, err := buf.WriteString(eol.EOL); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}
	return cntwrite, nil
}

func (t *Table) writeRecordHorBorder(data map[string]interface{}, w io.Writer) (int, error) {
	var cntwrite int
	cntCols := len(t.Columns)	
	if n, err := w.Write([]byte(Borders[t.border][BKLeftToRight]+Borders[t.border][BKHorizontal])); err == nil {
			cntwrite += n
	} else {
		return -1, err
	}
	
	for num, c := range t.Columns {
	    val, mok := data[c.Name]
		if !mok || val == nil {
			val = ""
		}	
		caption := fmt.Sprintf(t.GetMaskFormat(c), val)		
	    for  i := 0; i <utf8.RuneCountInString(caption)-1; i++ {		
		  if n, err := w.Write([]byte(Borders[t.border][BKHorizontal])); err == nil {
		
		cntwrite += n
	} else {
		return -1, err
	}
		}
		
		
		
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKBottomCross]
		} else {			
			delim = Borders[t.border][BKHorizontal]+Borders[t.border][BKRightToLeft]		
		}
		if n, err := w.Write([]byte(delim)); err == nil {
			cntwrite += n
		} else {
			return -1, err
		}			
		
	}
	
	if n, err := w.Write([]byte("\n")); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}
	return cntwrite, nil
}


func (t *Table) writeData(buf *bufio.Writer) (int, error) {
	var step int
	if t.dataget != nil {
		for {
			step=0;
			ok, data := t.dataget()
			if !ok {
				break
			}
			step=step+1;
			if (step>1)&&t.CloseEachColumn {						
			  if _, err := t.writeRecordHorBorder(data, buf); err != nil {
				return -1, err					
			    }
			}		
			
			
			if _, err := t.writeRecord(data, buf); err != nil {
				return -1, err
			}
		}
	} else if t.CountData() != 0 {
		var ll int =len(t.Data)
		for ii, data := range t.Data {
			if _, err := t.writeRecord(data, buf); err != nil {
				return -1, err
			}
			if t.CloseEachColumn {			
			if ii < ll-1{
			   if _, err := t.writeRecordHorBorder(data, buf); err != nil {		
				return -1, err
			      }		
			   }
			}
		}
	}
	return buf.Buffered(), buf.Flush()
}

func (t *Table) autoWidth() error {
	//each column
	var wa columns.Columns
	colsvisbile := t.Columns.ColumnsVisible()
	for i := range colsvisbile {
		if colsvisbile[i].IsAutoSize() || t.autoSize > 0 {
			colsvisbile[i].MaxLen = len(colsvisbile[i].Caption)

			wa.Add(colsvisbile[i])
		}
	}
	if len(wa) == 0 {
		return nil
	}
	for _, data := range t.Data {
		for i := range wa {
			curval := fmt.Sprintf("%v", data[wa[i].Name])
			curlen := utf8.RuneCountInString(curval)
			if curlen > wa[i].MaxLen {
				wa[i].MaxLen = curlen
			}
		}
	}
	//autosize table
	if t.autoSize > 0 {
		termwidth := t.autoSize - utf8.RuneCountInString(Borders[t.border][BKVertical])*colsvisbile.Len() - utf8.RuneCountInString(Borders[t.border][BKVerticalBorder])*2
		nowwidths := make([]int, colsvisbile.Len())
		allcolswidth := 0
		for i := range colsvisbile {
			if colsvisbile[i].MaxLen > colsvisbile[i].Width || colsvisbile[i].IsAutoSize() {
				nowwidths[i] = colsvisbile[i].MaxLen
			} else {
				nowwidths[i] = colsvisbile[i].Width
			}
			allcolswidth += nowwidths[i]
		}
		//todo: allcolswidth - borders
		twAll := 0
		for i := range colsvisbile {
			colsvisbile[i].MaxLen = int(math.Trunc(float64(termwidth) * (float64(nowwidths[i]) / float64(allcolswidth))))
			twAll += colsvisbile[i].MaxLen
		}
		i := 0
		//distrib mod
		for {
			if twAll >= termwidth || twAll <= 0 {
				break
			}
			if i+1 >= colsvisbile.Len() {
				i = 0
			}
			colsvisbile[i].MaxLen = colsvisbile[i].MaxLen + 1

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
	buf := bufio.NewWriter(w)
	cols := t.Columns.ColumnsVisible()
	if cols.Len() == 0 {
		return 0, nil
	}
	if err := t.autoWidth(); err != nil {
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
		panic(err)
	}
	return buf.String()
}
