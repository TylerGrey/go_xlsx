package utils

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
	"sort"
	"strconv"
	"time"
)

const (
	DefaultSheetName = "Sheet1"
	SignOmit         = "-"
	TagColumnName    = "col_name"
	TagColumnOrder   = "col_order"
)

type ColumnType struct {
	Field string
	Title string
}

type target struct {
	field string
	tag   string
	order int
}

type Excel struct {
	sheet      string
	startRow   int
	datasource interface{}
	columns    []ColumnType
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

	swWrapper := streamWriter{sw}
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
	defer timeTrack(now, "Flush")
	if err = swWrapper.Flush(); err != nil {
		return nil, err
	}

	return f, nil
}

func (e *Excel) makeDefaultColumns() (columnTypes []ColumnType) {
	now := time.Now()
	defer timeTrack(now, "makeDefaultColumns")

	targets := e.targetFieldsAndTags()

	for i := 0; i < len(targets); i++ {
		title := targets[i].field
		if len(targets[i].tag) > 0 {
			title = targets[i].tag
		}

		columnTypes = append(columnTypes, ColumnType{
			Field: targets[i].field,
			Title: title,
		})
	}

	return
}

func (e *Excel) targetFieldsAndTags() (targets []target) {
	now := time.Now()
	defer timeTrack(now, "targetFieldsAndTags")

	elem := reflect.TypeOf(e.datasource).Elem()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := field.Tag.Get(TagColumnName)

		order, err := strconv.Atoi(field.Tag.Get(TagColumnOrder))
		if err != nil {
			order = elem.NumField()
		}

		if tag != SignOmit {
			targets = append(targets, target{
				field: field.Name,
				tag:   tag,
				order: order,
			})
		}
	}

	sort.Slice(targets, func(i, j int) bool {
		return targets[i].order < targets[j].order
	})

	return
}

type streamWriter struct {
	*excelize.StreamWriter
}

func (f *streamWriter) setHeader(e *Excel) error {
	now := time.Now()
	defer timeTrack(now, "setHeader")

	var headers []interface{}
	for _, column := range e.columns {
		headers = append(headers, column.Title)
	}

	if err := f.SetRow(fmt.Sprintf("A%d", e.startRow), headers); err != nil {
		return err
	}

	return nil
}

func (f *streamWriter) setBody(e *Excel) error {
	now := time.Now()
	defer timeTrack(now, "setBody")

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

	return nil
}

func timeTrack(start time.Time, method string) {
	elapsed := time.Since(start)
	fmt.Printf("%s %s\n", method, elapsed)
}
