package main
import (
	"os"
	"github.com/arteev/fmttab"
)

func main() {	
	tab := fmttab.New("Table without border:",fmttab.BorderNone,nil)
	tab.AddColumn("Key",15,fmttab.AlignLeft).
		AddColumn("Value",15,fmttab.AlignRight)
	tab.Data = []map[string]interface{}{
		{
			"Key":"Key One",
			"Value": "Value One",
		},
		{
			"Key":"Key Two",
			"Value": "Value Two",
		},
	}
	tab.WriteTo(os.Stdout)
}