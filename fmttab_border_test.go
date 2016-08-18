package fmttab

import (
	"testing"
	"fmt"
	"github.com/arteev/fmttab/eol"
)

func TestBorderWithout(t *testing.T) {
	tab := New("Table", BorderNone, nil)
	tab.AddColumn("Column1", WidthAuto, AlignLeft)
	tab.AddColumn("C2", WidthAuto, AlignLeft)
	org := fmt.Sprintf("Table%[1]sColumn1    C2  %[1]s1234567890 test%[1]s", eol.EOL)

	tab.AppendData(map[string]interface{}{
		"Column1": "1234567890",
		"C2":      "test",
	})
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
}



func TestBorderHorizontalOne(t *testing.T) {	
	
	org := fmt.Sprintf("Table%[1]s╔═══════╗%[1]s║Column1║%[1]s╟───────╢%[1]s║test   ║%[1]s╚═══════╝%[1]s", eol.EOL)    
	
	tab := New("Table",BorderDouble,nil)
        tab.AddColumn("Column1",WidthAuto,AlignLeft)
	tab.CloseEachColumn=true;
		
	tab.AppendData(map[string]interface{} {
            "Column1": "test",            
        })
	
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
	}
func TestBorderHorizontalMulti(t *testing.T) {	
	
	org := fmt.Sprintf("Table%[1]s╔═══════╤═══════╤═══════╗%[1]s║Column1│Column2│Column3║%[1]s╟───────┼───────┼───────╢%[1]s║test   │       │       ║ %[1]s╟───────┼───────┼───────╢%[1]s║test2  │       │       ║%[1]s╚═══════╧═══════╧═══════╝%[1]s", eol.EOL)    
	
	tab := New("Table",BorderDouble,nil)
        tab.AddColumn("Column1",WidthAuto,AlignLeft)
		tab.AddColumn("Column2",WidthAuto,AlignLeft)
		tab.AddColumn("Column3",WidthAuto,AlignLeft)
	tab.CloseEachColumn=true;
	 tab.AppendData(map[string]interface{} {
            "Column1": "test",
            
        })
		 tab.AppendData(map[string]interface{} {
            "Column1": "test2",
            
        })
	
	res := tab.String()
	if org != res {
		t.Errorf("Excepted \n%q, got:\n%q", org, res)
	}
	}	
	
func TestBorderHorizontalTWo(t *testing.T) {		
    org2 := fmt.Sprintf("Table%[1]s╔═══════╗%[1]s║Column1║%[1]s╟───────╢%[1]s║test   ║%[1]s╟───────╢%[1]s║test2  ║%[1]s╚═══════╝%[1]s", eol.EOL)
	tab := New("Table",BorderDouble,nil)
        tab.AddColumn("Column1",WidthAuto,AlignLeft)
	tab.CloseEachColumn=true;
		
	tab.AppendData(map[string]interface{} {
            "Column1": "test",            
        })
	
   tab.AppendData(map[string]interface{} {
            "Column1": "test2",
            
        })
	res2 := tab.String()	
	if org2 != res2 {
		t.Errorf("Excepted \n%q, got:\n%q", org2, res2)
	}	
	}
func TestBorderHorizontalWithout(t *testing.T) {	
	org3 := fmt.Sprintf("Table%[1]s╔═══════╗%[1]s║Column1║%[1]s╟───────╢%[1]s╚═══════╝%[1]s", eol.EOL)
	tab2 := New("Table",BorderDouble,nil)
        tab2.AddColumn("Column1",WidthAuto,AlignLeft)
	tab2.CloseEachColumn=true;
	res3 := tab2.String()	
	if org3 != res3 {
		t.Errorf("Excepted \n%q, got:\n%q", org3, res3)
	}		
}

func TestReplaceBorder(t *testing.T) {
	tab := New("Table", BorderNone, nil)
	if tab.GetBorder() != BorderNone {
		t.Errorf("Excepted %v, got:%v", BorderNone, tab.border)
	}
	tab.SetBorder(BorderDouble)
	if tab.GetBorder() != BorderDouble {
		t.Errorf("Excepted %v, got:%v", BorderDouble, tab.border)
	}
}
