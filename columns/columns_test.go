package columns

import (
	"testing"
)

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
	columns[0].Visible = false
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
		t.Errorf("Excepted IsAutoSize=false, got %q", col1.IsAutoSize())
	}
	col1.Width = WidthAuto
	if !col1.IsAutoSize() {
		t.Errorf("Excepted IsAutoSize=true, got %q", col1.IsAutoSize())
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

	if col1 != columns[0] {
		t.Errorf("Excepted %v,got %v", col1, columns[0])
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

}
