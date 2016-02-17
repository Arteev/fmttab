package main

import (
	"os"
	"strings"

	"github.com/arteev/fmttab"
	"github.com/nsf/termbox-go"
)

func main() {
	tab := fmttab.New("Environments", fmttab.BorderDouble, nil)
	tab.AddColumn("ENV", fmttab.WidthAuto, fmttab.AlignLeft).
		AddColumn("VALUE", 20, fmttab.AlignLeft)
	for _, env := range os.Environ() {
		keyval := strings.Split(env, "=")
		tab.AppendData(map[string]interface{}{
			"ENV":   keyval[0],
			"VALUE": keyval[1],
		})
	}
	tab.WriteTo(os.Stdout)

	//Autofit
	err := termbox.Init()	
	if err != nil {
		panic(err)
	}
	termwidth, _ := termbox.Size()
	termbox.Close()    
	tab.AutoSize(true, termwidth)
	tab.WriteTo(os.Stdout)

}
