package main

import (
	"os"

	"github.com/arteev/fmttab"
)

func main() {
	borderMy := fmttab.Border(100)
	fmttab.Borders[borderMy] = map[fmttab.BorderKind]string{

		fmttab.BKVertical:         "|",
		fmttab.BKHorizontal:       "-",
		fmttab.BKHorizontalBorder: "_",
	}
	tab := fmttab.New("Custom table:", borderMy, nil)
	tab.AddColumn("Key", 15, fmttab.AlignLeft).
		AddColumn("Value", 15, fmttab.AlignRight)
	tab.Data = []map[string]interface{}{
		{
			"Key":   "Key One",
			"Value": "Value One",
		},
		{
			"Key":   "Key Two",
			"Value": "Value Two",
		},
	}
	tab.WriteTo(os.Stdout)
}
