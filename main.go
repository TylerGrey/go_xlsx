package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"net/http"
	"time"
	"xlsx_test/utils"
)

func main() {
	now := time.Now()
	datasource := []utils.TestType{
		{
			A: "A1",
			B: "B1",
			C: "C1",
			D: "D1",
		},
		{
			A: "A2",
			B: "B2",
			C: "C2",
			D: "D2",
		},
		{
			A: "A3",
			B: "B3",
			C: "C3",
			D: "D3",
		},
	}

	columns := []utils.ColumnType{}
	for i := 0; i < 3; i++ {
		columns = append(columns, utils.ColumnType{
			ID:      "",
			Title:   "TITLE",
			ColSpan: 0,
			Render:  nil,
		})
	}

	f, _ := utils.
		NewExcel().
		SetSheet("Test Sheet").
		SetStartRow(3).
		SetDataSource(datasource).
		//SetColumns(columns).
		Render()

	if err := f.SaveAs("Test.xlsx"); err != nil {
		fmt.Println(err)
	}

	fmt.Println(time.Since(now))

	//http.HandleFunc("/xlsx", downloadExcel)
	//http.ListenAndServe(":3000", nil)
}

func PrepareAndReturnExcel() *excelize.File {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Username")
	f.SetCellValue("Sheet1", "A2", "CC")
	f.SetCellValue("Sheet1", "B1", "Location")
	f.SetCellValue("Sheet1", "B2", "DD")
	f.SetCellValue("Sheet1", "C1", "Occupation")
	f.SetCellValue("Sheet1", "C2", "DD")
	return f
}

func downloadExcel(w http.ResponseWriter, r *http.Request) {
	// Get the Excel file with the user input data
	file := PrepareAndReturnExcel()

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=userInputData.xlsx")
	w.Header().Set("File-Name", "userInputData.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
}
