package columns

import (
	"errors"
)

const (
	//WidthAuto auto sizing of width column
	WidthAuto = 0
	//AlignLeft align text along the left edge
	AlignLeft = Align(false)
	//AlignRight align text along the right edge
	AlignRight = Align(true)
)

//Errors
var (
	ErrorAlreadyExists = errors.New("Column already exists")
)

//A Align text alignment in column of the table
type Align bool

//A Column type of table columns
type Column struct {
	MaxLen  int
	Caption string
	Name    string
	Width   int
	Aling   Align
	Visible bool
}

//A Columns array of the columns
type Columns struct {
	columns []*Column
}

//IsAutoSize returns auto sizing of width column
func (t Column) IsAutoSize() bool {
	return t.Width == WidthAuto
}

//Len returns count columns
func (c *Columns) Len() int {
	if c == nil {
		return 0
	}
	return len(c.columns)
}

//FindByName returns columns by name if exists or nil
func (c *Columns) FindByName(name string) *Column {
	for i, col := range c.columns {
		if col.Name == name {
			return c.columns[i]
		}
	}
	return nil
}

//NewColumn append new column in list with check by name of column
func (c *Columns) NewColumn(name, caption string, width int, aling Align) (*Column, error) {
	if c.FindByName(name) != nil {
		return nil, ErrorAlreadyExists
	}
	column := &Column{
		Name:    name,
		Caption: caption,
		Width:   width,
		Aling:   aling,
		Visible: true,
	}
	c.columns = append(c.columns, column)
	return column, nil
}

//Add append column with check exists
func (c *Columns) Add(col *Column) error {
	for i := range c.columns {
		if c.columns[i] == col {
			return ErrorAlreadyExists
		}
	}
	if c.FindByName(col.Name) != nil {
		return ErrorAlreadyExists
	}

	c.columns = append(c.columns, col)
	return nil
}

//ColumnsVisible returns count visible columns
func (c *Columns) ColumnsVisible() (res Columns) {
	for _, col := range c.columns {
		if col.Visible {
			res.columns = append(res.columns, col)
		}
	}
	return
}

//Columns gets array columns
func (c *Columns) Columns() []Column {
	var result []Column
	for _, col := range c.columns {
		result = append(result, *col)
	}
	return result
}

func (c *Columns) Get(index int) *Column {
	return c.columns[index]
}

func (c *Columns) Visit(f func(c *Column) error) error {
	for i := range c.columns {
		err := f(c.columns[i])
		if err != nil {
			return err
		}
	}
	return nil
}
