package fmttab
import (
	"bytes"

	"io"

	"strings"
	"fmt"
	"strconv"
)

type Border int
type borderKind int
type Align bool
var (
	BORDER_THIN = Border(1)
	BORDER_DOUBLE = Border(2)
	ALIGN_LEFT = Align(false)
	ALIGN_RIGHT = Align(true)
)
const (
	bkBetween  borderKind = iota
	bkLeftTop
	bkRighttop
	bkRightBottom
	bkLeftBottom
	bkLeftToRight
	bkRightToLeft
	bkTopToBottom
	bkBottomToTop
	bkBottomCross
	bkHorizontal
	bkVertical
	bkHorizontalBorder
	bkVerticalBorder
)

const (
	trimend = ".."
)

var borders = map[Border]map[borderKind]string{
	BORDER_THIN:map[borderKind]string{
		bkBetween : "\u2502",
		bkLeftTop : "\u250c",
		bkRighttop: "\u2510",
		bkRightBottom: "\u2518",
		bkLeftBottom: "\u2514",
		bkLeftToRight: "\u251c",
		bkRightToLeft: "\u2524",
		bkTopToBottom: "\u252c",
		bkBottomToTop: "\u2534",
		bkBottomCross: "\u253c",
		bkHorizontal: "\u2500",
		bkVertical: "\u2502",
		bkHorizontalBorder:"\u2500",
		bkVerticalBorder:"\u2502",
	},
	BORDER_DOUBLE : map[borderKind]string{
		bkBetween : "\u2502",
		bkLeftTop : "\u2554",
		bkRighttop: "\u2557",
		bkRightBottom: "\u255d",
		bkLeftBottom: "\u255a",
		bkLeftToRight: "\u255f",
		bkRightToLeft: "\u2562",
		bkTopToBottom: "\u2564",
		bkBottomToTop: "\u2567",
		bkBottomCross: "\u253c",
		bkHorizontal:  "\u2500",
		bkVertical: "\u2502",
		bkHorizontalBorder:"\u2550",
		bkVerticalBorder:"\u2551",
	},

}

type DataGetter func() (bool, map[string]interface{})

type column struct {
	name  string
	width int
	aling Align

}

type Table struct {
	cols    []*column
	caption string
	border  Border
	dataget DataGetter
}


func TrimEnds(val,end string, max int) string {
	l := len(val)
	if l<=max {
		return val
	}
	lend := len(end)
	if lend>=max {
		return end[:max]
	}
	return val[:(max-lend)] + end
}
func (c *column) GetMaskFormat() string {

	if c.aling== ALIGN_LEFT {
		return 	"%-" + strconv.Itoa(c.width) + "v"
	}
	return 	"%" + strconv.Itoa(c.width) + "v"
}

func (t *Table) AddColumn(name string, width int, aling Align) *Table {
	t.cols = append(t.cols, &column{
		name:name,
		width:width,
		aling:aling,
	})
	return t
}

func (t *Table) writeHeader(w io.Writer) (int,error) {
	var cntwrite int

	if n, err := w.Write([]byte(borders[t.border][bkLeftTop]));err != nil {
		return -1,err
	} else {
		cntwrite+=n
	}
	cntCols := len(t.cols)
	for num, c := range t.cols {
		if n, err := w.Write([]byte( strings.Repeat(borders[t.border][bkHorizontalBorder], c.width)));err != nil {
			return -1,err
		} else {
			cntwrite+=n
		}
		var delim string
		if num < cntCols - 1 {
			delim = borders[t.border][bkTopToBottom]
		} else {
			delim = borders[t.border][bkRighttop] + "\n"
		}

		if n, err := w.Write([]byte(delim));err != nil {
			return -1,err
		} else {
			cntwrite+=n
		}
	}

	if n, err := w.Write([]byte(borders[t.border][bkVerticalBorder]));err != nil {
		return -1,err
	} else {
		cntwrite+=n
	}
	for num, c := range t.cols {
		caption := fmt.Sprintf( c.GetMaskFormat(), c.name)

		if n, err := w.Write([]byte( TrimEnds(caption,trimend,c.width) ));err != nil {
			return -1,err
		} else {
			cntwrite+=n
		}

		var delim string
		if num < cntCols - 1 {
			delim = borders[t.border][bkVertical]
		} else {
			delim = borders[t.border][bkVerticalBorder] + ""
		}


		if n, err := w.Write([]byte(delim));err != nil {
			return -1, err
		} else {
			cntwrite+=n
		}
	}

	if n, err := w.Write([]byte("\n")); err != nil {
		return -1,err
	}else {
		cntwrite+=n
	}
	if n, err := w.Write([]byte(borders[t.border][bkLeftToRight]));err != nil {
		return 1, err
	}else {
		cntwrite+=n
	}
	for num, c := range t.cols {

		if n, err := w.Write([]byte( strings.Repeat(borders[t.border][bkHorizontal], c.width)));err != nil {
			return -1,err
		}else {
			cntwrite+=n
		}
		var delim string
		if num < cntCols - 1 {
			delim = borders[t.border][bkBottomCross]
		} else {
			delim = borders[t.border][bkRightToLeft] + "\n"
		}

		if n, err := w.Write([]byte(delim));err != nil {
			return -1,err
		}else {
			cntwrite+=n
		}
	}
	return cntwrite,nil
}

func (t *Table) writeBottomBorder(w io.Writer) (int,error) {
	var cntwrite int
	if n, err := w.Write([]byte(borders[t.border][bkLeftBottom])); err != nil {
		return -1,err
	} else {
		cntwrite+=n
	}
	cntCols := len(t.cols)
	for num, c := range t.cols {
		if n, err := w.Write([]byte( strings.Repeat(borders[t.border][bkHorizontalBorder], c.width))); err != nil {
			return -1,err
		} else {
			cntwrite+=n
		}
		var delim string
		if num < cntCols - 1 {
			delim = borders[t.border][bkBottomToTop]
		} else {
			delim = borders[t.border][bkRightBottom] + "\n"
		}
		if n, err := w.Write([]byte(delim)); err != nil {
			return -1,err
		} else {
			cntwrite+=n
		}
	}
	return cntwrite,nil
}


func (t *Table) writeData(w io.Writer) (int,error) {
	var cntwrite int
	if t.dataget != nil {
		cntCols := len(t.cols)
		for {
			ok, data := t.dataget()
			if !ok {
				break
			}
			if n, err := w.Write([]byte(borders[t.border][bkVerticalBorder])); err != nil {
				return -1,err
			} else {
				cntwrite += n
			}
			for num, c := range t.cols {
				val, mok := data[c.name];
				if !mok || val==nil {
					val = ""
				}
				caption := fmt.Sprintf(c.GetMaskFormat(), val)
				if n, err := w.Write([]byte(  TrimEnds(caption,trimend,c.width) )); err != nil {
					return -1,err
				} else {
					cntwrite += n
				}
				var delim string
				if num < cntCols - 1 {
					delim = borders[t.border][bkVertical]
				} else {
					delim = borders[t.border][bkVerticalBorder] + ""
				}
				if n, err := w.Write([]byte(delim)); err != nil {
					return -1,err
				} else {
					cntwrite += n
				}
			}
			if n, err := w.Write([]byte("\n")); err != nil {
				return -1,err
			} else {
				cntwrite += n
			}
		}
	}
	return cntwrite,nil
}

func New(caption string, border Border, datagetter DataGetter) *Table {
	return &Table{
		caption: caption,
		border:border,
		dataget:datagetter,
	}
}

func (t *Table) String() string {
	var buf bytes.Buffer
	t.WriteTo(&buf)
	return buf.String()
}

func (t *Table) WriteTo(w io.Writer)(int64,error){
	if (len(t.cols)==0) {
		return 0,nil
	}
	var cntwrite int64
	if n,err:=t.writeHeader(w); err!=nil {
		return -1,err
	}else {
		cntwrite += int64(n)
	}
	if n,err:=t.writeData(w);err!=nil {
		return -1,err
	} else {
		cntwrite += int64(n)
	}
	if n,err:=t.writeBottomBorder(w); err!=nil {
		return -1,err
	} else {
		cntwrite += int64(n)
	}
	return cntwrite,nil
}

