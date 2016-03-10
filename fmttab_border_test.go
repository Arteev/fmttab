package fmttab

import (	
	"testing"
    
)

func TestBorderWithout(t *testing.T) {
	tab := New("Table", BorderNone, nil)
	tab.AddColumn("Column1", WidthAuto, AlignLeft)	
    tab.AddColumn("C2", WidthAuto, AlignLeft)
    org := "Table\nColumn1    C2  \n1234567890 test\n"
	tab.AppendData(map[string]interface{}{
		"Column1": "1234567890",
        "C2": "test",
	})
	res := tab.String()    
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
}

func TestReplaceBorder(t *testing.T) {
    tab := New("Table", BorderNone, nil)
    if tab.GetBorder()!=BorderNone {
        t.Errorf("Excepted %v, got:%v", BorderNone, tab.border)
    }
    tab.SetBorder(BorderDouble)
    if tab.GetBorder()!=BorderDouble {
        t.Errorf("Excepted %v, got:%v", BorderDouble, tab.border)
    }
}