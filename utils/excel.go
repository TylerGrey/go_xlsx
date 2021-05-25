package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"reflect"
	"sort"
	"strconv"
)

const (
	DefaultSheetName   = "Sheet1"
	DefaultStyleHeader = "{\"border\":[{\"type\":\"left\",\"color\":\"000000\",\"style\":1},{\"type\":\"right\",\"color\":\"000000\",\"style\":1},{\"type\":\"top\",\"color\":\"000000\",\"style\":1},{\"type\":\"bottom\",\"color\":\"000000\",\"style\":1}],\"font\":{\"bold\":true},\"fill\":{\"type\":\"pattern\",\"color\":[\"cccccc\"],\"pattern\":1},\"alignment\":{\"horizontal\":\"center\",\"vertical\":\"center\"}}"
	DefaultStyleBody   = "{\"border\":[{\"type\":\"left\",\"color\":\"000000\",\"style\":1},{\"type\":\"right\",\"color\":\"000000\",\"style\":1},{\"type\":\"top\",\"color\":\"000000\",\"style\":1},{\"type\":\"bottom\",\"color\":\"000000\",\"style\":1}],\"alignment\":{\"vertical\":\"center\"}}"

	SignOmit = "-"

	TagColumnName   = "col_name"
	TagColumnOrder  = "col_order"
	TagColumnRender = "col_render"
)

// ColumnType 컬럼 정보
type ColumnType struct {
	// Field 렌더링 될 구조체 필드 명
	Field string
	// Name 컬럼 명
	Name string
	// MergeColumn AutoMerging 시 병합 여부
	MergeColumn bool
	// Children AutoMerging 시 sub column
	Children []ColumnType
	// Render 렌더링 될 데이터 변경을 위한 callback
	Render func(v interface{}) interface{}
}

// target 컬럼 자동 생성 시, 렌더링 될 컬럼 정보
type target struct {
	field  string
	tag    string
	order  int
	render func(v interface{}) interface{}
}

// Excel ...
type Excel struct {
	sheet      string
	startRow   int
	datasource interface{}
	columns    []ColumnType
	autoMerge  bool
}

// NewExcel Default 구조체 생성
func NewExcel() *Excel {
	return &Excel{
		sheet:     DefaultSheetName,
		startRow:  1,
		autoMerge: false,
	}
}

// SetSheet 시트 명 설정
func (e *Excel) SetSheet(sheet string) *Excel {
	e.sheet = sheet
	return e
}

// SetStartRow 렌더링 될 시작 행 위치 설정
func (e *Excel) SetStartRow(startRow int) *Excel {
	if startRow < 1 {
		startRow = 1
	}
	e.startRow = startRow
	return e
}

// SetDataSource 렌더링 할 데이터 설정. datasource 는 slice 만 사용 가능
func (e *Excel) SetDataSource(datasource interface{}) *Excel {
	e.datasource = datasource
	return e
}

// SetColumns 렌더링 할 컬럼 설정. 설정 안할 시, struct tag 바탕으로 default 컬럼 생성 됨
func (e *Excel) SetColumns(columns []ColumnType) *Excel {
	e.columns = columns
	return e
}

// SetAutoMerge 자동 셀 병합이 필요할 경우 설정
func (e *Excel) SetAutoMerge(autoMerge bool) *Excel {
	e.autoMerge = autoMerge
	return e
}

// Render 전달 받은 정보를 바탕으로 Excelize file 생성
func (e *Excel) Render() (*excelize.File, error) {
	if e.datasource == nil || reflect.TypeOf(e.datasource).Kind() != reflect.Slice {
		return nil, errors.New("invalid datasource")
	}

	f := excelize.NewFile()
	f.SetSheetName(DefaultSheetName, e.sheet)

	if e.columns == nil {
		e.columns = e.makeDefaultColumns()
	}

	// 2.3.2 버전 StreamWriter 셀 병합 미지원으로 코드 분리
	if e.autoMerge {
		file := file{f}

		leafColumns, headerHeight, err := file.drawHeader(e)
		if err != nil {
			return nil, err
		}

		if err = file.drawBody(leafColumns, headerHeight, e); err != nil {
			return nil, err
		}
	} else {
		sw, err := f.NewStreamWriter(e.sheet)
		if err != nil {
			return nil, err
		}
		swWrapper := streamWriter{sw}

		if err = swWrapper.drawHeader(e); err != nil {
			return nil, err
		}

		if err = swWrapper.drawBody(e); err != nil {
			return nil, err
		}

		if err = swWrapper.Flush(); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// calcWidth Sub column 까지 포함한 컬럼의 총 길이를 구하는 함수
func (e *Excel) calcWidth(c ColumnType) int {
	width := 0

	if c.Children == nil {
		return 1
	}

	for _, child := range c.Children {
		width += e.calcWidth(child)
	}

	return width
}

// makeDefaultColumns 데이터 모델의 tag를 바탕으로 default column 생성
func (e *Excel) makeDefaultColumns() (columnTypes []ColumnType) {
	targets := e.targetFieldsAndTags()

	for i := 0; i < len(targets); i++ {
		title := targets[i].field
		if len(targets[i].tag) > 0 {
			title = targets[i].tag
		}

		columnTypes = append(columnTypes, ColumnType{
			Field:  targets[i].field,
			Name:   title,
			Render: targets[i].render,
		})
	}

	return
}

// newStyleHeader default header 스타일 반환
func (e *Excel) newStyleHeader() excelize.Style {
	var style excelize.Style
	json.Unmarshal([]byte(DefaultStyleHeader), &style)

	return style
}

// newStyleBody default body 스타일 반환
func (e *Excel) newStyleBody() excelize.Style {
	var style excelize.Style
	json.Unmarshal([]byte(DefaultStyleBody), &style)

	return style
}

// targetFieldsAndTags 컬럼 자동 생성 시, tag 기반으로 렌더링 될 컬럼 정보 생성
func (e *Excel) targetFieldsAndTags() (targets []target) {
	elem := reflect.TypeOf(e.datasource).Elem()

	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	// 데이터 모델의 필드에 대한 loop
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := field.Tag.Get(TagColumnName)

		order, err := strconv.Atoi(field.Tag.Get(TagColumnOrder))
		if err != nil {
			order = elem.NumField()
		}

		renderTag := field.Tag.Get(TagColumnRender)
		var render func(v interface{}) interface{}
		if len(renderTag) > 0 {
			render = func(v interface{}) interface{} {
				method := reflect.New(elem).MethodByName(renderTag)
				out := method.Call([]reflect.Value{reflect.ValueOf(v)})
				return out[0]
			}
		}

		if tag != SignOmit {
			targets = append(targets, target{
				field:  field.Name,
				tag:    tag,
				order:  order,
				render: render,
			})
		}
	}

	// sorting
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].order < targets[j].order
	})

	return
}

// file Automerging 인 경우, file 로 생성. (StreamWriter 해당 버전은 셀 병합 미지원)
type file struct {
	*excelize.File
}

// drawHeader 헤더 렌더링
func (f *file) drawHeader(e *Excel) ([]ColumnType, int, error) {
	var leafColumnsSize int
	for _, column := range e.columns {
		leafColumnsSize += e.calcWidth(column)
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

				if column.Name == prevCellValue {
					err := f.MergeCell(e.sheet, prevCell, startCell)
					if err != nil {
						return nil, 0, err
					}

					newColumns = append(newColumns, column)
					continue
				}
			}

			if column.Children != nil {
				columnWidth := e.calcWidth(column) - 1
				forcePosition += columnWidth

				mergeColumnName, _ := excelize.ColumnNumberToName(startColumnNumber + columnWidth)

				if err := f.SetCellValue(e.sheet, startCell, column.Name); err != nil {
					return nil, 0, err
				}
				// horizontal merging
				if err := f.MergeCell(e.sheet, startCell, fmt.Sprintf("%s%d", mergeColumnName, e.startRow+headerHeight)); err != nil {
					return nil, 0, err
				}

				newColumns = append(newColumns, column.Children...)
			} else {
				if err := f.SetCellValue(e.sheet, startCell, column.Name); err != nil {
					return nil, 0, err
				}
				newColumns = append(newColumns, column)
			}
		}

		if len(leafColumns) >= leafColumnsSize {
			break
		}

		leafColumns = newColumns
		headerHeight++
	}

	// header style 적용
	style := e.newStyleHeader()
	styleId, err := f.NewStyle(&style)
	if err != nil {
		return nil, 0, err
	}

	startName, _ := excelize.ColumnNumberToName(1)
	startCell := fmt.Sprintf("%s%d", startName, e.startRow)
	endName, _ := excelize.ColumnNumberToName(leafColumnsSize)
	endCell := fmt.Sprintf("%s%d", endName, e.startRow+headerHeight)

	if err := f.SetCellStyle(e.sheet, startCell, endCell, styleId); err != nil {
		return nil, 0, err
	}

	return leafColumns, headerHeight, nil
}

// drawBody Body 렌더링
func (f *file) drawBody(leafColumns []ColumnType, headerHeight int, e *Excel) error {
	style := e.newStyleBody()
	styleId, _ := f.NewStyle(&style)

	valueOf := reflect.ValueOf(e.datasource)
	for i := 0; i < valueOf.Len(); i++ {
		isMergedRow := true
		currentRow := e.startRow + headerHeight + i + 1

		valueOfData := valueOf.Index(i)
		if valueOfData.Kind() == reflect.Ptr {
			valueOfData = valueOfData.Elem()
		}

		for j, column := range leafColumns {
			columnName, _ := excelize.ColumnNumberToName(j + 1)
			currentCell := fmt.Sprintf("%s%d", columnName, currentRow)

			value := valueOfData.FieldByName(column.Field).Interface()
			if column.Render != nil {
				value = column.Render(value)
			}

			if i > 0 && column.MergeColumn {
				prevCell := fmt.Sprintf("%s%d", columnName, currentRow-1)
				prevCellValue, _ := f.GetCellValue(e.sheet, prevCell)

				if prevCellValue == value && isMergedRow {
					if err := f.MergeCell(e.sheet, prevCell, currentCell); err != nil {
						return err
					}
				} else {
					isMergedRow = false
					if err := f.SetCellValue(e.sheet, currentCell, value); err != nil {
						return err
					}
				}
			} else {
				if err := f.SetCellValue(e.sheet, currentCell, value); err != nil {
					return err
				}
			}

			if err := f.SetCellStyle(e.sheet, currentCell, currentCell, styleId); err != nil {
				return err
			}
		}
	}

	return nil
}

// streamWriter 일반적인 경우, 성능을 위해 스트림 이용
type streamWriter struct {
	*excelize.StreamWriter
}

// drawHeader Header 렌더링
func (f *streamWriter) drawHeader(e *Excel) error {
	style := e.newStyleHeader()
	styleId, _ := f.File.NewStyle(&style)

	var headers []interface{}
	for _, column := range e.columns {
		headers = append(headers, excelize.Cell{
			StyleID: styleId,
			Value:   column.Name,
		})
	}

	if err := f.SetRow(fmt.Sprintf("A%d", e.startRow), headers); err != nil {
		return err
	}

	return nil
}

// drawBody Body 렌더링
func (f *streamWriter) drawBody(e *Excel) error {
	style := e.newStyleBody()
	styleId, _ := f.File.NewStyle(&style)

	valueOf := reflect.ValueOf(e.datasource)

	for i := 0; i < valueOf.Len(); i++ {
		var rows []interface{}

		for j := 0; j < len(e.columns); j++ {
			valueOfData := valueOf.Index(i)

			if valueOfData.Kind() == reflect.Ptr {
				valueOfData = valueOfData.Elem()
			}

			currentValue := valueOfData.FieldByName(e.columns[j].Field).Interface()
			if e.columns[j].Render != nil {
				currentValue = e.columns[j].Render(currentValue)
			}

			rows = append(rows, excelize.Cell{
				StyleID: styleId,
				Value:   currentValue,
			})
		}

		if err := f.SetRow(fmt.Sprintf("A%d", (i+1)+e.startRow), rows); err != nil {
			return err
		}
	}

	return nil
}
