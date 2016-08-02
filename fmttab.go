package fmttab

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"math"

	"github.com/arteev/fmttab/columns"
	"github.com/arteev/fmttab/eol"
	"github.com/arteev/tern"
)

//A Border of table
type Border int

//A BorderKind type of element on the border of the table
type BorderKind int

const (
	//WidthAuto auto sizing of width column
	WidthAuto = columns.WidthAuto
	//BorderNone table without borders
	BorderNone = Border(0)
	//BorderThin table with a thin border
	BorderThin = Border(1)
	//BorderDouble table with a double border
	BorderDouble = Border(2)

	//BorderDouble table with a simple border
	BorderSimple = Border(3)
	//AlignLeft align text along the left edge
	AlignLeft = columns.AlignLeft
	//AlignRight align text along the right edge
	AlignRight = columns.AlignRight
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

//Trimend - end of line after trimming
var Trimend = ".."

//Borders predefined border types
var Borders = map[Border]map[BorderKind]string{
	BorderNone: map[BorderKind]string{
		BKVertical: " ",
	},
	BorderSimple: map[BorderKind]string{
		BKBottomCross: "+",
		BKHorizontal:  "-",
		BKVertical:    "|",
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

//A Table is the repository for the columns, the data that are used for printing the table
type Table struct {
	dataget       DataGetter
	border        Border
	caption       string
	autoSize      int
	Columns       columns.Columns
	Data          []map[string]interface{}
	VisibleHeader bool
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
func (t *Table) GetMaskFormat(c *columns.Column) string {
	if c.Aling == AlignLeft {
		return "%-" + strconv.Itoa(t.getWidth(c)) + "v"
	}
	return "%" + strconv.Itoa(t.getWidth(c)) + "v"

}

//must be calculated before call
func (t *Table) getWidth(c *columns.Column) int {
	if c.IsAutoSize() || t.autoSize > 0 {
		return c.MaxLen
	}
	return c.Width

}

//AddColumn adds a column to the table
func (t *Table) AddColumn(name string, width int, aling columns.Align) *Table {
	_, err := t.Columns.NewColumn(name, name, width, aling)
	if err != nil {
		panic(err)
	}
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

//AutoSize fit columns
func (t *Table) AutoSize(enabled bool, destWidth int) {
	if enabled {
		t.autoSize = destWidth
	} else {
		t.autoSize = 0
	}
}

//CountData the amount of data in the table
func (t *Table) CountData() int {
	return len(t.Data)
}

//SetBorder - set  type of border table
func (t *Table) SetBorder(b Border) {
	t.border = b
}

//GetBorder - get current border
func (t Table) GetBorder() Border {
	return t.border
}

func (t *Table) writeHeader(w io.Writer) (int, error) {
	var cntwrite int
	dataout := ""
	if t.caption != "" {
		dataout += t.caption + eol.EOL
	}
	dataout += Borders[t.border][BKLeftTop]
	colv := t.Columns.ColumnsVisible()
	cntCols := tern.Op(colv == nil, 0, colv.Len()).(int)
	olddataout := dataout
	for num, c := range colv {
		dataout += strings.Repeat(Borders[t.border][BKHorizontalBorder], t.getWidth(c))
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKTopToBottom]
		} else {
			delim = Borders[t.border][BKRighttop]
			if olddataout != dataout {
				delim += eol.EOL
			}
		}
		dataout += delim
	}

	if t.VisibleHeader {
		dataout += Borders[t.border][BKVerticalBorder]
		for num, c := range colv {
			caption := fmt.Sprintf(t.GetMaskFormat(c), c.Caption)
			dataout += trimEnds(caption, Trimend, t.getWidth(c))
			var delim string
			if num < cntCols-1 {
				delim = Borders[t.border][BKVertical]
			} else {
				delim = Borders[t.border][BKVerticalBorder]
			}
			dataout += delim
		}
		dataout += eol.EOL + Borders[t.border][BKLeftToRight]
		dataout += t.writeBorderTopButtomData(BKHorizontal, BKBottomCross, BKRightToLeft)
	}

	if n, err := w.Write([]byte(dataout)); err == nil {
		cntwrite += n
	} else {
		return -1, err
	}
	return cntwrite, nil
}

func (t *Table) writeBorderTopButtomData(hr, vbwnCol, vright BorderKind) (data string) {
	colv := t.Columns.ColumnsVisible()
	cntCols := tern.Op(colv == nil, 0, colv.Len()).(int)
	empty := true
	for num, c := range colv {
		s := strings.Repeat(Borders[t.border][hr], t.getWidth(c))
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
				delim += eol.EOL
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
	colv := t.Columns.ColumnsVisible()
	cntCols := tern.Op(colv == nil, 0, colv.Len()).(int)
	if n, err := w.Write([]byte(Borders[t.border][BKVerticalBorder])); err == nil {
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
		if n, err := w.Write([]byte(trimEnds(caption, Trimend, t.getWidth(c)))); err == nil {
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
	if n, err := w.Write([]byte(eol.EOL)); err == nil {
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
		caption:       caption,
		border:        border,
		dataget:       datagetter,
		VisibleHeader: true,
	}
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
	cols := t.Columns.ColumnsVisible()
	if cols.Len() == 0 {
		return 0, nil
	}
	if err := t.autoWidth(); err != nil {
		return 0, err
	}
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
