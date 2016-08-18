package fmttab

import (
	"io/ioutil"
	"testing"

	"math/rand"
)

func makeTable() *Table {
	tab := New("Table", BorderThin, nil)
	tab.AddColumn("Column1", 8, AlignLeft)
	tab.AddColumn("Column2", 16, AlignLeft)
	tab.CloseEachColumn = true
	for n := 0; n < 100; n++ {
		tab.AppendData(map[string]interface{}{
			"Column1": n * rand.Int(),
			"Column2": "data",
		})
	}
	return tab
}

func BenchmarkWriteTo(b *testing.B) {
	tab := makeTable()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := tab.WriteTo(ioutil.Discard)
		if err != nil {
			b.Error(err)
		}
	}

}

func BenchmarkString(b *testing.B) {
	tab := makeTable()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := tab.WriteTo(ioutil.Discard)
		if err != nil {
			b.Error(err)
		}
	}
}
