package fmttab

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

//A Border of table
type Border int

//A BorderKind type of element on the border of the table
type BorderKind int

//A Align text alignment in column of the table
type Align bool

const (
	//WidthAuto auto sizing of width column
	WidthAuto = 0
	//BorderNone table without borders
	BorderNone = Border(0)
	//BorderThin table with a thin border
	BorderThin = Border(1)
	//BorderDouble table with a double border
	BorderDouble = Border(2)
	//AlignLeft align text along the left edge
	AlignLeft = Align(false)
	//AlignRight align text along the right edge
	AlignRight = Align(true)
)

//The concrete type of the object on the border of the table
const (
	BKLeftTop BorderKind = iota
	BKRighttop
	BKRightBottom
	BKLeftBottom
	BKLeftToRight
	BKRightToLeft
	BKTopToBottom
	BKBottomToTop
	BKBottomCross
	BKHorizontal
	BKVertical
	BKHorizontalBorder
	BKVerticalBorder
)

const (
	trimend = ".."
)

//Borders predefined border types
var Borders = map[Border]map[BorderKind]string{
	BorderNone: map[BorderKind]string{
		BKVertical: "|",
	},
	BorderThin: map[BorderKind]string{
		BKLeftTop:          "\u250c",
		BKRighttop:         "\u2510",
		BKRightBottom:      "\u2518",
		BKLeftBottom:       "\u2514",
		BKLeftToRight:      "\u251c",
		BKRightToLeft:      "\u2524",
		BKTopToBottom:      "\u252c",
		BKBottomToTop:      "\u2534",
		BKBottomCross:      "\u253c",
		BKHorizontal:       "\u2500",
		BKVertical:         "\u2502",
		BKHorizontalBorder: "\u2500",
		BKVerticalBorder:   "\u2502",
	},
	BorderDouble: map[BorderKind]string{
		BKLeftTop:          "\u2554",
		BKRighttop:         "\u2557",
		BKRightBottom:      "\u255d",
		BKLeftBottom:       "\u255a",
		BKLeftToRight:      "\u255f",
		BKRightToLeft:      "\u2562",
		BKTopToBottom:      "\u2564",
		BKBottomToTop:      "\u2567",
		BKBottomCross:      "\u253c",
		BKHorizontal:       "\u2500",
		BKVertical:         "\u2502",
		BKHorizontalBorder: "\u2550",
		BKVerticalBorder:   "\u2551",
	},
}

//A DataGetter functional type for table data
type DataGetter func() (bool, map[string]interface{})

//A Column type of table columns
type Column struct {
	maxLen int
	Name   string
	Width  int
	Aling  Align
}

//A Table is the repository for the columns, the data that are used for printing the table
type Table struct {
	Columns []*Column
	caption string
	border  Border
	Data    []map[string]interface{}
	dataget DataGetter
}

// A trimEnds supplements the text with special characters by limiting the length of the text column width
func trimEnds(val, end string, max int) string {
	l := utf8.RuneCountInString(val)
	if l <= max {
		return val
	}
	lend := utf8.RuneCountInString(end)
	if lend >= max {
		return end[:max]
	}
	return string([]rune(val)[:(max-lend)]) + end
}

//GetMaskFormat returns a pattern string for formatting text in table column alignment
func (c *Column) GetMaskFormat() string {
	if c.Aling == AlignLeft {
		return "%-" + strconv.Itoa(c.getWidth()) + "v"
	}
	return "%" + strconv.Itoa(c.getWidth()) + "v"
}

//must be calculated before call
func (c *Column) getWidth() int {
	if c.Width == WidthAuto {
		return c.maxLen
	}
	return c.Width
}

//AddColumn adds a column to the table
func (t *Table) AddColumn(name string, width int, aling Align) *Table {
	//TODO: check dublicate
	t.Columns = append(t.Columns, &Column{
		Name:  name,
		Width: width,
		Aling: aling,
	})
	return t
}

//AppendData adds the data to the table
func (t *Table) AppendData(rec map[string]interface{}) *Table {
	t.Data = append(t.Data, rec)
	return t
}

//ClearData removes data from a table
func (t *Table) ClearData() *Table {
	t.Data = nil
	return t
}

//CountData the amount of data in the table
func (t *Table) CountData() int {
	return len(t.Data)
}

func (t *Table) writeHeader(w io.Writer) (int, error) {
	var cntwrite int
	dataout := ""
	if t.caption != "" {
		dataout += t.caption + "\n"
	}
	dataout += Borders[t.border][BKLeftTop]
	cntCols := len(t.Columns)
	for num, c := range t.Columns {
		dataout += strings.Repeat(Borders[t.border][BKHorizontalBorder], c.getWidth())
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKTopToBottom]
		} else {
			delim = Borders[t.border][BKRighttop] + "\n"
		}
		dataout += delim
	}
	dataout += Borders[t.border][BKVerticalBorder]
	for num, c := range t.Columns {
		caption := fmt.Sprintf(c.GetMaskFormat(), c.Name)
		dataout += trimEnds(caption, trimend, c.getWidth())
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKVertical]
		} else {
			delim = Borders[t.border][BKVerticalBorder]
		}
		dataout += delim
	}
	dataout += "\n" + Borders[t.border][BKLeftToRight]
	dataout += t.writeBorderTopButtomData(BKHorizontal, BKBottomCross, BKRightToLeft)
	if n, err := w.Write([]byte(dataout)); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}
	return cntwrite, nil
}

func (t *Table) writeBorderTopButtomData(hr, vbwnCol, vright BorderKind) (data string) {
	cntCols := len(t.Columns)
	empty := true
	for num, c := range t.Columns {
		s := strings.Repeat(Borders[t.border][hr], c.getWidth())
		if len(s) > 0 {
			empty = false
		}
		data += s
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][vbwnCol]
		} else {
			delim = Borders[t.border][vright]
			if !empty {
				delim += "\n"
			}
		}
		if len(delim) > 0 {
			empty = false
		}
		data += delim
	}
	return data
}

func (t *Table) writeBottomBorder(w io.Writer) (int, error) {
	var cntwrite int
	data := Borders[t.border][BKLeftBottom] + t.writeBorderTopButtomData(BKHorizontalBorder, BKBottomToTop, BKRightBottom)
	if n, err := w.Write([]byte(data)); err == nil {
		cntwrite = n
	} else {
		return -1, err
	}
	return cntwrite, nil
}

func (t *Table) writeRecord(data map[string]interface{}, w io.Writer) (int, error) {
	var cntwrite int
	cntCols := len(t.Columns)
	if n, err := w.Write([]byte(Borders[t.border][BKVerticalBorder])); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}
	for num, c := range t.Columns {
		val, mok := data[c.Name]
		if !mok || val == nil {
			val = ""
		}
		caption := fmt.Sprintf(c.GetMaskFormat(), val)
		if n, err := w.Write([]byte(trimEnds(caption, trimend, c.getWidth()))); err == nil {
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

func (t *Table) writeData(w io.Writer) (int, error) {
	var cntwrite int
	if t.dataget != nil {
		for {
			ok, data := t.dataget()
			if !ok {
				break
			}
			if n, err := t.writeRecord(data, w); err == nil {
				cntwrite += n
			} else {
				return -1, err
			}
		}
	} else if t.CountData() != 0 {
		for _, data := range t.Data {
			if n, err := t.writeRecord(data, w); err == nil {
				cntwrite += n
			} else {
				return -1, err
			}
		}
	}
	return cntwrite, nil
}

//New creates a Table object. DataGetter can be nil
func New(caption string, border Border, datagetter DataGetter) *Table {
	return &Table{
		caption: caption,
		border:  border,
		dataget: datagetter,
	}
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

func (t *Table) autoWidth() {
	var wa []*Column
	for i := range t.Columns {
		if t.Columns[i].Width == WidthAuto {
			t.Columns[i].maxLen = len(t.Columns[i].Name)
			wa = append(wa, t.Columns[i])
		}
	}
	if len(wa) == 0 {
		return
	}
	for _, data := range t.Data {
		for i := range wa {
			curval := fmt.Sprintf("%v", data[wa[i].Name])
			curlen := len(curval)
			if curlen > wa[i].maxLen {
				wa[i].maxLen = curlen
			}
		}
	}	
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	if len(t.Columns) == 0 {
		return 0, nil
	}
	t.autoWidth()
	var cntwrite int64
	if n, err := t.writeHeader(w); err == nil {
		cntwrite += int64(n)
	} else {
		return -1, err
	}
	if n, err := t.writeData(w); err == nil {
		cntwrite += int64(n)
	} else {
		return -1, err
	}
	if n, err := t.writeBottomBorder(w); err == nil {
		cntwrite += int64(n)
	} else {
		return -1, err
	}
	return cntwrite, nil
}
