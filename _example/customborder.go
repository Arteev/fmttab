package main
import (
	"os"
	"github.com/arteev/fmttab"
)

func main() {
	BORDER_MY := fmttab.Border(100)
	fmttab.Borders[BORDER_MY] = map[fmttab.BorderKind]string{

		fmttab.BKVertical: "|",
		fmttab.BKHorizontal : "-",
		fmttab.BKHorizontalBorder : "_",
	}
	tab := fmttab.New("Custom table:",BORDER_MY,nil)
	tab.AddColumn("Key",15,fmttab.ALIGN_LEFT).
		AddColumn("Value",15,fmttab.ALIGN_RIGHT)
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