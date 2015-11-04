package main
import (
	"os"
	"strings"
	"github.com/arteev/fmttab"
)
func main() {
	tab := fmttab.New("Environments:",fmttab.BORDER_NONE,nil)
	tab.Columns = []*fmttab.Column{
		{
			Name:"ENV",Width:25,Aling:fmttab.ALIGN_LEFT,
		},
		{
			Name:"VALUE",Width:25,Aling:fmttab.ALIGN_LEFT,
		},
	}
	for _,env:=range os.Environ() {
		keyval := strings.Split(env,"=")
		tab.Data = append(tab.Data,map[string]interface{} {
			"ENV": keyval[0],
			"VALUE" : keyval[1],
		})
	}
	tab.WriteTo(os.Stdout)
}