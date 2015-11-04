package main
import (
	"os"
	"strings"
	"github.com/arteev/fmttab"
)
func main() {
	tab := fmttab.New("Environments",fmttab.BORDER_DOUBLE,nil)
	tab.AddColumn("ENV",25,fmttab.ALIGN_LEFT).
		AddColumn("VALUE",25,fmttab.ALIGN_LEFT)
	for _,env:=range os.Environ() {
		keyval := strings.Split(env,"=")
		tab.AppendData(map[string]interface{} {
			"ENV": keyval[0],
			"VALUE" : keyval[1],
		})
	}
	tab.WriteTo(os.Stdout)
}