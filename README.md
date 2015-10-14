# fmttab

[![Build Status](https://travis-ci.org/Arteev/fmttab.svg?branch=master)](https://travis-ci.org/Arteev/fmttab)
[![Coverage Status](https://coveralls.io/repos/Arteev/fmttab/badge.svg?branch=master&service=github)](https://coveralls.io/github/Arteev/fmttab?branch=master)
Description
-----------

Package Golang to display information table pseudographics

Installation
------------

This package can be installed with the go get command:

    go get github.com/Arteev/fmttab

Documentation
-------------
Example 1:

```go
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
```

Output:
```
╔═════════════════════════╤═════════════════════════╗
║ENV                      │VALUE                    ║
╟─────────────────────────┼─────────────────────────╢
║PAPERSIZE                │a4                       ║
║UPSTART_SESSION          │unix:abstract            ║
║GNOME_DESKTOP_SESSION_ID │this-is-deprecated       ║
║GDMSESSION               │ubuntu                   ║
║IM_CONFIG_PHASE          │1                        ║
║COMPIZ_CONFIG_PROFILE    │ubuntu                   ║
║LANG                     │ru_RU.UTF-8              ║
╚═════════════════════════╧═════════════════════════╝
```

Example 2:
```go
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
```

Output:
```
╔══════════════════════════════╤══════════╤════════════════════╤══════╗
║Name                          │      Size│Time                │Dir   ║
╟──────────────────────────────┼──────────┼────────────────────┼──────╢
║example                       │      4096│2015-10-13 16:29:10 │true  ║
║console.go                    │       860│2015-10-13 16:22:10 │false ║
║example                       │   2669536│2015-10-13 16:29:10 │false ║
╚══════════════════════════════╧══════════╧════════════════════╧══════╝
```


License
-------

  MIT

Author
------

Arteev Aleksey
