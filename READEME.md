# EXCEL

```go
type Excel struct {
	Datasource []interface{}
	Columns []ColumnType
}

type ColumnType struct {
	ID  string
	Title string
	ColSpan int32
	Render *func()
}
```