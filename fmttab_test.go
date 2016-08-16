package fmttab

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/arteev/fmttab/eol"
)

func TestCallAndOut(t *testing.T) {
	const name = "<THIS IS NAME>"
	const val = 123
	arr := []map[string]interface{}{
		map[string]interface{}{
			"ID":   val,
			"NAME": name,
		},
	}
	l := len(arr)
	cur := 0
	tbl := New("", BorderDouble, func() (bool, map[string]interface{}) {
		if cur >= l {
			return false, nil
		}
		r := arr[cur]
		cur++
		return true, r
	})
	tbl.AddColumn("ID", 10, AlignLeft).
		AddColumn("NAME", 40, AlignRight)

	sout := tbl.String()
	if cur != 1 {
		t.Error("Excepted call datagetter")
	}

	if strings.IndexAny(sout, name) == -1 {
		t.Errorf("Excepted in output %q", name)
	}
	if strings.IndexAny(sout, strconv.Itoa(val)) == -1 {
		t.Errorf("Excepted in output %v", val)
	}

	tbl2 := New("", BorderDouble, nil)
	tbl2.AddColumn("ID", 10, AlignLeft).
		AddColumn("NAME", 40, AlignRight)

	for _, item := range arr {
		tbl2.AppendData(item)
	}
	sout2 := tbl2.String()
	if sout2 != sout {
		t.Errorf("Excepted (AppendData) %q, got %q", sout, sout2)
	}

}

func TestTrim(t *testing.T) {
	test := map[struct {
		val string
		end string
		max int
	}]string{
		{
			"testing", "...", 6,
		}: "tes...",
		{
			"testing", "...", 7,
		}: "testing",
		{
			"testing", "...", 2,
		}: "..",
		{
			"testing", ">", 5,
		}: "test>",
	}
	for key, pair := range test {
		r := trimEnds(key.val, key.end, key.max)
		if r != pair {
			t.Errorf("Excepted %q, got %q", pair, r)
		}
	}
}

func TestString(t *testing.T) {
	tab := New("Table", BorderThin, nil)
	tab.AddColumn("Column1", 8, AlignLeft)
	org := fmt.Sprintf("Table%[1]s┌────────┐%[1]s│Column1 │%[1]s├────────┤%[1]s└────────┘%[1]s", eol.EOL)
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}

	var buf bytes.Buffer

	if cntwrite, err := tab.WriteTo(&buf); err != nil {
		t.Error(err)
	} else {
		if int64(len(res)) != cntwrite {
			t.Errorf("Excepted count write:%d, got: %d", len(res), cntwrite)
		}
	}

	var buf2 bytes.Buffer
	tab2 := New("", BorderThin, nil)
	if cntwrite, err := tab2.WriteTo(&buf2); err != nil {
		t.Error(err)
	} else {
		if cntwrite != 0 {
			t.Errorf("Excepted count write:-1, got: %d", cntwrite)
		}
	}

}

func TestNationalChars(t *testing.T) {
	tab := New("Table", BorderThin, nil)
	tab.AddColumn("Column1", 16, AlignLeft)
	org := fmt.Sprintf("Table%[1]s┌────────────────┐%[1]s│Column1         │%[1]s├────────────────┤%[1]s│Русский         │%[1]s└────────────────┘%[1]s", eol.EOL)
	tab.AppendData(map[string]interface{}{
		"Column1": "Русский",
	})
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
}

func TestData(t *testing.T) {
	tab := New("", BorderThin, nil)
	if tab.CountData() != 0 {
		t.Errorf("Excepted 0, got:%d", tab.CountData())
	}
	data1 := map[string]interface{}{
		"COL1": "value1",
		"COL2": "value2",
	}
	tab.AppendData(data1)
	if tab.CountData() != 1 {
		t.Errorf("Excepted 1, got:%d", tab.CountData())
	}
	tab.ClearData()
	if tab.CountData() != 0 {
		t.Errorf("Excepted 0, got:%d", tab.CountData())
	}
}

func TestAutoWidthColumns(t *testing.T) {
	tab := New("Table", BorderThin, nil)
	tab.AddColumn("Column1", WidthAuto, AlignLeft)
	org := fmt.Sprintf("Table%[1]s┌──────────┐%[1]s│Column1   │%[1]s├──────────┤%[1]s│1234567890│%[1]s└──────────┘%[1]s", eol.EOL)
	tab.AppendData(map[string]interface{}{
		"Column1": "1234567890",
	})
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
}

func TestAutoSize(t *testing.T) {
	tab := New("Table", BorderThin, nil)
	tab.AddColumn("Column1", WidthAuto, AlignLeft)
	tab.AutoSize(false, 0)
	orgNormal := fmt.Sprintf("Table%[1]s┌──────────┐%[1]s│Column1   │%[1]s├──────────┤%[1]s│1234567890│%[1]s└──────────┘%[1]s", eol.EOL)
	orgAutoSize10 := fmt.Sprintf("Table%[1]s┌───────┐%[1]s│Column1│%[1]s├───────┤%[1]s│12345..│%[1]s└───────┘%[1]s", eol.EOL)
	orgAutoSize16 := fmt.Sprintf("Table%[1]s┌─────────────┐%[1]s│Column1      │%[1]s├─────────────┤%[1]s│1234567890   │%[1]s└─────────────┘%[1]s", eol.EOL)
	tab.AppendData(map[string]interface{}{
		"Column1": "1234567890",
	})
	res := tab.String()
	if orgNormal != res {
		t.Errorf("Excepted \n%q, got:\n%q", orgNormal, res)
	}

	tab.AutoSize(true, 16)
	res = tab.String()
	if orgAutoSize16 != res {
		t.Errorf("Excepted \n%q, got:\n%q", orgAutoSize16, res)
	}

	tab.AutoSize(true, 10)
	res = tab.String()
	if orgAutoSize10 != res {
		t.Errorf("Excepted \n%q, got:\n%q", orgAutoSize10, res)
	}
}

func TestHideColumn(t *testing.T) {
	tab := New("Table", BorderThin, nil)
	tab.AddColumn("Column1", 8, AlignLeft)
	tab.AddColumn("Column2", 8, AlignLeft)
	tab.Columns[tab.Columns.Len()-1].Visible = false
	org := fmt.Sprintf("Table%[1]s┌────────┐%[1]s│Column1 │%[1]s├────────┤%[1]s└────────┘%[1]s", eol.EOL)
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
	data1 := map[string]interface{}{
		"Column1": "value1",
		"Column2": "value2",
	}
	tab.AppendData(data1)
	org2 := fmt.Sprintf("Table%[1]s┌────────┐%[1]s│Column1 │%[1]s├────────┤%[1]s│value1  │%[1]s└────────┘%[1]s", eol.EOL)
	res = tab.String()
	if org2 != res {
		t.Errorf("Excepted \n%q, got:\n%q", org2, res)
	}

	tab.Columns[0].Visible = false
	res = tab.String()
	if "" != res {
		t.Errorf("Excepted \n%q, got:\n%q", "", res)
	}
}

func TestWithoutHeader(t *testing.T) {
	tab := New("Table", BorderThin, nil)
	tab.VisibleHeader = false
	tab.AddColumn("Column1", WidthAuto, AlignLeft)
	org := fmt.Sprintf("Table%[1]s┌──────────┐%[1]s│1234567890│%[1]s└──────────┘%[1]s", eol.EOL)
	tab.AppendData(map[string]interface{}{
		"Column1": "1234567890",
	})
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
}
