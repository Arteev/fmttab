package fmttab

import (
	"bytes"

	"io"

	"fmt"
	"strconv"
	"strings"
)

type Border int
type BorderKind int
type Align bool

var (
	BORDER_NONE   = Border(0)
	BORDER_THIN   = Border(1)
	BORDER_DOUBLE = Border(2)
	ALIGN_LEFT    = Align(false)
	ALIGN_RIGHT   = Align(true)
)

const (
	BKLeftTop = iota
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

var Borders = map[Border]map[BorderKind]string{
	BORDER_NONE: map[BorderKind]string{
		BKVertical: "|",
	},

	BORDER_THIN: map[BorderKind]string{
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
	BORDER_DOUBLE: map[BorderKind]string{
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

type DataGetter func() (bool, map[string]interface{})

type Column struct {
	Name  string
	Width int
	Aling Align
}

type Table struct {
	Columns []*Column
	caption string
	border  Border
	Data    []map[string]interface{}
	dataget DataGetter
}

func TrimEnds(val, end string, max int) string {
	l := len(val)
	if l <= max {
		return val
	}
	lend := len(end)
	if lend >= max {
		return end[:max]
	}
	return val[:(max-lend)] + end
}

func (c *Column) GetMaskFormat() string {

	if c.Aling == ALIGN_LEFT {
		return "%-" + strconv.Itoa(c.Width) + "v"
	}
	return "%" + strconv.Itoa(c.Width) + "v"
}

func (t *Table) AddColumn(name string, width int, aling Align) *Table {
	t.Columns = append(t.Columns, &Column{
		Name:  name,
		Width: width,
		Aling: aling,
	})
	return t
}

func (t *Table) AppendData(rec map[string]interface{}) *Table {
	t.Data = append(t.Data, rec)
	return t
}

func (t *Table) ClearData() *Table {
	t.Data = nil
	return t
}

func (t *Table) CountData() int {
	return len(t.Data)
}

func (t *Table) writeHeader(w io.Writer) (int, error) {
	var cntwrite int

	if t.caption != "" {
		if n, err := w.Write([]byte(t.caption + "\n")); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}
	}

	if n, err := w.Write([]byte(Borders[t.border][BKLeftTop])); err != nil {
		return -1, err
	} else {
		cntwrite += n
	}
	cntCols := len(t.Columns)
	for num, c := range t.Columns {
		if n, err := w.Write([]byte(strings.Repeat(Borders[t.border][BKHorizontalBorder], c.Width))); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKTopToBottom]
		} else {
			delim = Borders[t.border][BKRighttop] + "\n"
		}

		if n, err := w.Write([]byte(delim)); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}
	}

	if n, err := w.Write([]byte(Borders[t.border][BKVerticalBorder])); err != nil {
		return -1, err
	} else {
		cntwrite += n
	}
	for num, c := range t.Columns {
		caption := fmt.Sprintf(c.GetMaskFormat(), c.Name)

		if n, err := w.Write([]byte(TrimEnds(caption, trimend, c.Width))); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}

		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKVertical]
		} else {
			delim = Borders[t.border][BKVerticalBorder] + ""
		}

		if n, err := w.Write([]byte(delim)); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}
	}

	if n, err := w.Write([]byte("\n")); err != nil {
		return -1, err
	} else {
		cntwrite += n
	}
	if n, err := w.Write([]byte(Borders[t.border][BKLeftToRight])); err != nil {
		return 1, err
	} else {
		cntwrite += n
	}

	emptyhdr := true
	for num, c := range t.Columns {
		if n, err := w.Write([]byte(strings.Repeat(Borders[t.border][BKHorizontal], c.Width))); err != nil {
			return -1, err
		} else {
			cntwrite += n
			if n > 0 {
				emptyhdr = false
			}
		}
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKBottomCross]
		} else {
			delim = Borders[t.border][BKRightToLeft]
			if !emptyhdr {
				delim += "\n"
			}
		}

		if n, err := w.Write([]byte(delim)); err != nil {
			return -1, err
		} else {
			cntwrite += n
			if n > 0 {
				emptyhdr = false
			}
		}
	}
	return cntwrite, nil
}

func (t *Table) writeBottomBorder(w io.Writer) (int, error) {
	var cntwrite int
	if n, err := w.Write([]byte(Borders[t.border][BKLeftBottom])); err != nil {
		return -1, err
	} else {
		cntwrite += n
	}
	cntCols := len(t.Columns)
	empty := true
	for num, c := range t.Columns {
		if n, err := w.Write([]byte(strings.Repeat(Borders[t.border][BKHorizontalBorder], c.Width))); err != nil {
			return -1, err
		} else {
			cntwrite += n
			if n > 0 {
				empty = false
			}
		}
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKBottomToTop]
		} else {
			delim = Borders[t.border][BKRightBottom]
			if !empty {
				delim += "\n"
			}
		}
		if n, err := w.Write([]byte(delim)); err != nil {
			return -1, err
		} else {
			cntwrite += n
			if n > 0 {
				empty = false
			}
		}
	}
	return cntwrite, nil
}

func (t *Table) writeRecord(data map[string]interface{}, w io.Writer) (int, error) {
	var cntwrite int
	cntCols := len(t.Columns)
	if n, err := w.Write([]byte(Borders[t.border][BKVerticalBorder])); err != nil {
		return -1, err
	} else {
		cntwrite += n
	}
	for num, c := range t.Columns {
		val, mok := data[c.Name]
		if !mok || val == nil {
			val = ""
		}
		caption := fmt.Sprintf(c.GetMaskFormat(), val)
		if n, err := w.Write([]byte(TrimEnds(caption, trimend, c.Width))); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}
		var delim string
		if num < cntCols-1 {
			delim = Borders[t.border][BKVertical]
		} else {
			delim = Borders[t.border][BKVerticalBorder] + ""
		}
		if n, err := w.Write([]byte(delim)); err != nil {
			return -1, err
		} else {
			cntwrite += n
		}
	}
	if n, err := w.Write([]byte("\n")); err != nil {
		return -1, err
	} else {
		cntwrite += n
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
			if n, err := t.writeRecord(data, w); err != nil {
				return -1, err
			} else {
				cntwrite += n
			}
		}
	} else if t.CountData() != 0 {
		for _, data := range t.Data {
			if n, err := t.writeRecord(data, w); err != nil {
				return -1, err
			} else {
				cntwrite += n
			}
		}
	}
	return cntwrite, nil
}

func New(caption string, border Border, datagetter DataGetter) *Table {
	return &Table{
		caption: caption,
		border:  border,
		dataget: datagetter,
	}
}

func (t *Table) String() string {
	var buf bytes.Buffer
	t.WriteTo(&buf)
	return buf.String()
}

func (t *Table) WriteTo(w io.Writer) (int64, error) {
	if len(t.Columns) == 0 {
		return 0, nil
	}
	var cntwrite int64
	if n, err := t.writeHeader(w); err != nil {
		return -1, err
	} else {
		cntwrite += int64(n)
	}
	if n, err := t.writeData(w); err != nil {
		return -1, err
	} else {
		cntwrite += int64(n)
	}
	if n, err := t.writeBottomBorder(w); err != nil {
		return -1, err
	} else {
		cntwrite += int64(n)
	}
	return cntwrite, nil
}
