package columns

import (
	"errors"
	"reflect"
	"testing"
)

func TestVisit(t *testing.T) {
	var columns Columns
	c1, err := columns.NewColumn("Col1", "Columns 1", 10, AlignLeft)
	if err != nil {
		t.Error(err)
	}
	c2, err := columns.NewColumn("Col2", "Columns 2", 10, AlignLeft)
	if err != nil {
		t.Error(err)
	}
	countVisit := 0
	err = columns.Visit(func(c *Column) error {
		if c1 != c && countVisit == 0 {
			t.Errorf("Expected %v , got %v", c1, c)
		}
		if c2 != c && countVisit == 1 {
			t.Errorf("Expected %v , got %v", c2, c)
		}
		countVisit++
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	if countVisit != 2 {
		t.Errorf("Expected visit %d,actual %d", 2, countVisit)
	}

	err = columns.Visit(func(c *Column) error {
		return errors.New("test")
	})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestColumnsVisible(t *testing.T) {
	var columns Columns
	if cls := columns.ColumnsVisible(); cls.Len() != 0 {
		t.Errorf("Excepted CountVisible %d, got %d", 0, cls.Len())
	}
	_, err := columns.NewColumn("Col1", "Columns 1", 10, AlignLeft)
	if err != nil {
		t.Error(err)
	}
	if cls := columns.ColumnsVisible(); cls.Len() != 1 {
		t.Errorf("Excepted CountVisible %d, got %d", 1, cls.Len())
	}
	_, err = columns.NewColumn("Col2", "Columns 2", 10, AlignLeft)
	if err != nil {
		t.Error(err)
	}
	if cls := columns.ColumnsVisible(); cls.Len() != 2 {
		t.Errorf("Excepted CountVisible %d, got %d", 2, cls.Len())
	}
	columns.columns[0].Visible = false
	if cls := columns.ColumnsVisible(); cls.Len() != 1 {
		t.Errorf("Excepted CountVisible %d, got %d", 1, cls.Len())
	}
}

func TestColumns(t *testing.T) {
	var (
		columns Columns
		pcoll   *Columns
	)
	if pcoll.Len() != 0 {
		t.Errorf("Excepted len %d, got %d", 0, pcoll.Len())
	}
	if columns.Len() != 0 {
		t.Errorf("Excepted len %d, got %d", 0, columns.Len())
	}
	col1, err := columns.NewColumn("Col1", "Columns 1", 10, AlignLeft)
	if err != nil {
		t.Error(err)
	}
	if col1.Name != "Col1" {
		t.Errorf("Excepted %q, got %q", "Col1", col1.Name)
	}
	if col1.IsAutoSize() {
		t.Errorf("Excepted IsAutoSize=false, got %t", col1.IsAutoSize())
	}
	col1.Width = WidthAuto
	if !col1.IsAutoSize() {
		t.Errorf("Excepted IsAutoSize=true, got %t", col1.IsAutoSize())
	}

	if columns.Len() != 1 {
		t.Errorf("Excepted len %d, got %d", 1, columns.Len())
	}

	colf := columns.FindByName("Col1")
	if colf != col1 {
		t.Errorf("Excepted %v, got %v", col1, colf)
	}
	colnot := columns.FindByName("?")
	if colnot != nil {
		t.Errorf("Excepted nil,got %v", colnot)
	}

	if col1 != columns.columns[0] {
		t.Errorf("Excepted %v,got %v", col1, columns.columns[0])
	}

	err = columns.Add(col1)
	if err != ErrorAlreadyExists {
		t.Errorf("Excepted %v, got %v", ErrorAlreadyExists, err)
	}
	col2 := &Column{
		Name: "Col2",
	}
	err = columns.Add(col2)
	if err != nil {
		t.Error(err)
	}

	_, err = columns.NewColumn("Col2", "", 0, AlignLeft)
	if err != ErrorAlreadyExists {
		t.Errorf("Excepted %v, got %v", ErrorAlreadyExists, err)
	}

	err = columns.Add(&Column{Name: col1.Name})
	if err != ErrorAlreadyExists {
		t.Errorf("Excepted %v, got %v", ErrorAlreadyExists, err)
	}

	got := columns.Get(0)
	if got != col1 {
		t.Errorf("Exptected %v,got %v", col1, got)
	}

}

func TestColumn(t *testing.T) {
	var columns Columns
	col, err := columns.NewColumn("Col1", "Columns 1", 10, AlignLeft)
	if err != nil {
		t.Fatal(err)
	}
	got := col.GetWidth()
	if got != 10 {
		t.Errorf("Expected %d,got %d", 10, got)
	}

	col.MaxLen = 15
	got = col.GetWidth()
	if got != col.MaxLen {
		t.Errorf("Expected %d,got %d", col.MaxLen, got)
	}

	mask := "%-15v"
	if got := col.GetMaskFormat(); got != mask {
		t.Errorf("Expected %s,got %s", mask, got)
	}

	col.Aling = AlignRight
	mask = "%15v"
	if got := col.GetMaskFormat(); got != mask {
		t.Errorf("Expected %s,got %s", mask, got)
	}
}

func TestGetColumns(t *testing.T) {
	var columns Columns
	_, err := columns.NewColumn("Col1", "Columns 1", 10, AlignLeft)
	if err != nil {
		t.Fatal(err)
	}
	_, err = columns.NewColumn("Col2", "Columns 2", 15, AlignLeft)
	if err != nil {
		t.Fatal(err)
	}
	col2 := columns.Columns()
	if len(col2) != len(columns.columns) {
		t.Fatalf("Expected count columns %d,got %d", len(columns.columns), len(col2))
	}
	for i := 0; i < len(col2); i++ {
		if !reflect.DeepEqual(col2[i], *columns.columns[i]) {
			t.Errorf("Expected %v,got %v", *columns.columns[i], col2[i])
		}
	}
}
