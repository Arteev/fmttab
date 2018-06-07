package main

import (
	"os"
	"strings"

	"github.com/arteev/fmttab"
	"github.com/arteev/fmttab/columns"
)

func main() {
	tab := fmttab.New("Environments:", fmttab.BorderNone, nil)
	tab.AddColumn("ENV", 25, fmttab.AlignLeft)
	tab.Columns.Add(&columns.Column{
		Name: "VALUE", Width: 25, Aling: fmttab.AlignLeft, Visible: true,
	})

	for _, env := range os.Environ() {
		keyval := strings.Split(env, "=")
		tab.Data = append(tab.Data, map[string]interface{}{
			"ENV":   keyval[0],
			"VALUE": keyval[1],
		})
	}
	tab.WriteTo(os.Stdout)
}
