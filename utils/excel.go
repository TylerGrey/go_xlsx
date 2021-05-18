package utils

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
)

const (
	DefaultSheetName = "Sheet1"
	SignOmit         = "-"
	TagColumnName    = "col_name"
)

type Excel struct {
	sheet      string
	startRow   int
	datasource interface{}
	columns    []ColumnType
}

type ColumnType struct {
	ID    string
	Title string
}

func NewExcel() *Excel {
	return &Excel{
		sheet:    DefaultSheetName,
		startRow: 1,
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

	f := excelize.NewFile()
	f.SetSheetName(DefaultSheetName, e.sheet)
	sw, err := f.NewStreamWriter(e.sheet)

	swWrapper := StreamWriter{sw}
	if err != nil {
		return nil, err
	}

	if e.columns == nil {
		e.columns = e.makeDefaultColumns()
	}
	if err := swWrapper.setHeader(e); err != nil {
		return nil, err
	}

	if err := swWrapper.setBody(e); err != nil {
		return nil, err
	}

	if err = swWrapper.Flush(); err != nil {
		return nil, err
	}

	return f, nil
}

func (e *Excel) makeDefaultColumns() (columnTypes []ColumnType) {
	fields, tags := e.targetFieldsAndTags()

	for i := 0; i < len(fields); i++ {
		title := fields[i]
		if len(tags[i]) > 0 {
			title = tags[i]
		}

		columnTypes = append(columnTypes, ColumnType{
			ID:    fields[i],
			Title: title,
		})
	}

	return
}

func (e *Excel) targetFieldsAndTags() (fields []string, tags []string) {
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

type StreamWriter struct {
	*excelize.StreamWriter
}

func (f *StreamWriter) setHeader(e *Excel) error {
	var headers []interface{}
	for _, column := range e.columns {
		headers = append(headers, column.Title)
	}

	if err := f.SetRow(fmt.Sprintf("A%d", e.startRow), headers); err != nil {
		return err
	}

	return nil
}

func (f *StreamWriter) setBody(e *Excel) error {
	fields, _ := e.targetFieldsAndTags()

	valueOf := reflect.ValueOf(e.datasource)
	for i := 0; i < valueOf.Len(); i++ {
		var rows []interface{}

		for j := 0; j < len(fields); j++ {
			rows = append(rows, valueOf.Index(i).FieldByName(fields[j]))
		}

		if err := f.SetRow(fmt.Sprintf("A%d", (i+1)+e.startRow), rows); err != nil {
			return err
		}
	}

	return nil
}
