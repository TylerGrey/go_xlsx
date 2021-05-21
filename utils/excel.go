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
	Field    string
	Title    string
	Children []ColumnType
	// Render value 컨버팅 용
	Render func(v interface{}) interface{}
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
	if err != nil {
		return nil, err
	}

	swWrapper := streamWriter{sw}

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

func getWidth(c ColumnType) int {
	width := 0

	if c.Children == nil {
		return 1
	}

	for _, child := range c.Children {
		width += getWidth(child)
	}

	return width
}

func (e *Excel) RenderAutoMerging() (*excelize.File, error) {
	if e.datasource == nil || reflect.TypeOf(e.datasource).Kind() != reflect.Slice {
		return nil, errors.New("invalid datasource")
	}

	f := excelize.NewFile()
	f.SetSheetName(DefaultSheetName, e.sheet)

	// header
	var leafColumnsSize int
	for _, column := range e.columns {
		leafColumnsSize += getWidth(column)
	}

	// leafColumns 최종 columns
	leafColumns := e.columns
	var headerHeight int

	for {
		var newColumns []ColumnType
		var forcePosition int

		for i, column := range leafColumns {
			startColumnNumber := i + forcePosition + 1
			startColumnName, _ := excelize.ColumnNumberToName(startColumnNumber)
			startCell := fmt.Sprintf("%s%d", startColumnName, e.startRow+headerHeight)

			// Header 2번째 라인부터 vertical merging
			if headerHeight > 0 {
				prevCell := fmt.Sprintf("%s%d", startColumnName, e.startRow+headerHeight-1)
				prevCellValue, _ := f.GetCellValue(e.sheet, prevCell)

				if column.Title == prevCellValue {
					err := f.MergeCell(e.sheet, prevCell, startCell)
					if err != nil {
						return nil, err
					}

					newColumns = append(newColumns, column)
					continue
				}
			}

			if column.Children != nil {
				columnWidth := getWidth(column) - 1
				forcePosition += columnWidth

				mergeColumnName, _ := excelize.ColumnNumberToName(startColumnNumber + columnWidth)

				err := f.SetCellValue(e.sheet, startCell, column.Title)
				if err != nil {
					return nil, err
				}
				// horizontal merging
				err = f.MergeCell(e.sheet, startCell, fmt.Sprintf("%s%d", mergeColumnName, e.startRow+headerHeight))
				if err != nil {
					return nil, err
				}

				newColumns = append(newColumns, column.Children...)
			} else {
				err := f.SetCellValue(e.sheet, startCell, column.Title)
				if err != nil {
					return nil, err
				}
				newColumns = append(newColumns, column)
			}
		}

		if len(leafColumns) == leafColumnsSize {
			break
		}

		leafColumns = newColumns
		headerHeight++
	}

	// body
	valueOf := reflect.ValueOf(e.datasource)
	for i := 0; i < valueOf.Len(); i++ {
		for j, column := range leafColumns {
			columnName, _ := excelize.ColumnNumberToName(j + 1)

			value := valueOf.Index(i).FieldByName(column.Field).Interface()
			if column.Render != nil {
				value = column.Render(value)
			}

			err := f.SetCellValue(e.sheet, fmt.Sprintf("%s%d", columnName, e.startRow+headerHeight+i+1), value)
			if err != nil {
				return nil, err
			}
		}
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
