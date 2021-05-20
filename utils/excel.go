package utils

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
	"time"
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
	Field string
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

	now := time.Now()
	if err = swWrapper.Flush(); err != nil {
		return nil, err
	}
	fmt.Println("Flush", time.Since(now))

	return f, nil
}

func (e *Excel) makeDefaultColumns() (columnTypes []ColumnType) {
	now := time.Now()
	fields, tags := e.targetFieldsAndTags()

	for i := 0; i < len(fields); i++ {
		title := fields[i]
		if len(tags[i]) > 0 {
			title = tags[i]
		}

		columnTypes = append(columnTypes, ColumnType{
			Field: fields[i],
			Title: title,
		})
	}

	fmt.Println("makeDefaultColumns", time.Since(now))
	return
}

func (e *Excel) targetFieldsAndTags() (fields []string, tags []string) {
	now := time.Now()
	// TODO ordering 적용
	elem := reflect.TypeOf(e.datasource).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := field.Tag.Get(TagColumnName)

		if tag != SignOmit {
			fields = append(fields, field.Name)
			tags = append(tags, tag)
		}
	}

	fmt.Println("targetFieldsAndTags", time.Since(now))
	return
}

type StreamWriter struct {
	*excelize.StreamWriter
}

func (f *StreamWriter) setHeader(e *Excel) error {
	now := time.Now()
	var headers []interface{}
	for _, column := range e.columns {
		headers = append(headers, column.Title)
	}

	if err := f.SetRow(fmt.Sprintf("A%d", e.startRow), headers); err != nil {
		return err
	}

	fmt.Println("setHeader", time.Since(now))
	return nil
}

func (f *StreamWriter) setBody(e *Excel) error {
	now := time.Now()
	valueOf := reflect.ValueOf(e.datasource)
	for i := 0; i < valueOf.Len(); i++ {
		var rows []interface{}

		for j := 0; j < len(e.columns); j++ {
			rows = append(rows, valueOf.Index(i).FieldByName(e.columns[j].Field))
		}

		if err := f.SetRow(fmt.Sprintf("A%d", (i+1)+e.startRow), rows); err != nil {
			return err
		}
	}

	fmt.Println("setBody", time.Since(now))
	return nil
}
