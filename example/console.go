package main

import (
	"github.com/Arteev/fmttab"
	"os"
	"path/filepath"
)

var files []map[string]interface{}

func walkpath(path string, f os.FileInfo, err error) error {
	files = append(files, map[string]interface{} {
		"Name": f.Name(),
		"Size": f.Size(),
		"Dir": f.IsDir(),
		"Time": f.ModTime().Format("2006-01-02 15:04:01"),
	})
	return nil
}

func main() {
	files=make([]map[string]interface{},0)
	root := filepath.Dir(os.Args[0])
	filepath.Walk(root,walkpath)
	i:=0
	lfiles:=len(files)

	tab := fmttab.New("Table",fmttab.BORDER_DOUBLE,func() (bool, map[string]interface{}) {
		if i>=lfiles {
			return false,nil
		}
		i++
		return true,files[i-1]
	})
	tab.AddColumn("Name",30,fmttab.ALIGN_LEFT).
		AddColumn("Size",10,fmttab.ALIGN_RIGHT).
		AddColumn("Time",20,fmttab.ALIGN_LEFT).
		AddColumn("Dir",6,fmttab.ALIGN_LEFT)

	tab.WriteTo(os.Stdout)
}
