package utils

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
)

const (
	DefaultSheetName = "Sheet1"

	SignOmit = "-"

	TagColumnName = "col_name"
)

type Excel struct {
	sheet      string
	startRow   int
	datasource interface{}
	columns    []ColumnType
	autoMerge  bool
}

type ColumnType struct {
	ID      string
	Title   string
	ColSpan int32
	Render  *func()
}

func NewExcel() *Excel {
	// TODO fill all fields
	return &Excel{
		sheet:     DefaultSheetName,
		startRow:  1,
		autoMerge: false,
	}
}

func (e *Excel) SetSheet(sheet string) *Excel {
	e.sheet = sheet
	return e
}

func (e *Excel) SetStartRow(startRow int) *Excel {
	e.startRow = startRow
	return e
}

func (e *Excel) SetDataSource(datasource interface{}) *Excel {
	e.datasource = datasource
	return e
}

func (e *Excel) SetColumns(columns []ColumnType) *Excel {
	e.columns = columns
	return e
}

func (e *Excel) Render() (*excelize.File, error) {
	if e.datasource == nil || reflect.TypeOf(e.datasource).Kind() != reflect.Slice {
		return nil, errors.New("invalid datasource")
	}

	f := File{excelize.NewFile()}
	f.SetSheetName(DefaultSheetName, e.sheet)

	if e.columns == nil {
		e.columns = e.makeDefaultColumns()
	}
	if err := f.setHeader(e); err != nil {
		return nil, err
	}

	if err := f.setBody(e); err != nil {
		return nil, err
	}

	return f.File, nil
}

func (e *Excel) makeDefaultColumns() (columnTypes []ColumnType) {
	columns, tags := e.printableColumns()

	for i := 0; i < len(columns); i++ {
		title := columns[i]
		if len(tags[i]) > 0 {
			title = tags[i]
		}

		columnTypes = append(columnTypes, ColumnType{
			ID:      columns[i],
			Title:   title,
			ColSpan: 0,
			Render:  nil,
		})
	}

	return
}

func (e *Excel) printableColumns() (fields []string, tags []string) {
	elem := reflect.TypeOf(e.datasource).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := field.Tag.Get(TagColumnName)

		if tag != SignOmit {
			fields = append(fields, field.Name)
			tags = append(tags, tag)
		}
	}

	return
}

type File struct {
	*excelize.File
}

func (f *File) setHeader(e *Excel) error {
	for i := 1; i <= len(e.columns); i++ {
		colName, _ := excelize.ColumnNumberToName(i)
		axis := fmt.Sprintf("%s%d", colName, e.startRow)

		if err := f.SetCellValue(e.sheet, axis, e.columns[i-1].Title); err != nil {
			return err
		}
	}

	return nil
}

func (f *File) setBody(e *Excel) error {
	fields, _ := e.printableColumns()

	valueOf := reflect.ValueOf(e.datasource)
	for i := 1; i <= valueOf.Len(); i++ {
		axisRow := (i + 1) + e.startRow

		for j := 0; j < len(fields); j++ {
			axisName, _ := excelize.ColumnNumberToName(j)
			fmt.Println(fmt.Sprintf("%s%d", axisName, axisRow))

			field := valueOf.Index(i - 1).FieldByName(fields[j])
			fmt.Println("field", field)
			if err := f.SetCellValue(e.sheet, fmt.Sprintf("%s%d", axisName, axisRow), field); err != nil {
				return err
			}
		}
	}

	return nil
}
