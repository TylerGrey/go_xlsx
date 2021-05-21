package main

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"xlsx_test/utils"
)

func main() {
	// File
	//f := PrepareAndReturnExcel()
	//
	//now := time.Now()
	//if err := f.SaveAs("Test.xlsx"); err != nil {
	//	return
	//}
	//fmt.Println("SaveAs", time.Since(now))

	// HTTP
	//http.HandleFunc("/xlsx", downloadExcel)
	//http.ListenAndServe(":3000", nil)

	// Uncommon case
	PrepareUncommonCase().SaveAs("merge.xlsx")
}

func PrepareAndReturnExcel() *excelize.File {
	now := time.Now()
	var datasource []utils.UserInfo
	for i := 0; i < 1_000; i++ {
		datasource = append(datasource, utils.UserInfo{
			ID:             strconv.Itoa(i),
			Name:           fmt.Sprintf("Name%d", i),
			MobileNumber:   fmt.Sprintf("010-%4d-%4d", rand.Intn(9999), rand.Intn(9999)),
			EmployeeRegNum: strconv.Itoa(i),
			TeamName:       fmt.Sprintf("Team %d", rand.Int()),
			CompanyEmail:   fmt.Sprintf("%d@gmail.com", rand.Int()),
		})
	}
	fmt.Println("Generate datasource", time.Since(now))

	f, _ := utils.
		NewExcel().
		SetSheet("Test Sheet").
		SetStartRow(2).
		SetDataSource(datasource).
		Render()

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

	now := time.Now()
	file.Write(w)
	fmt.Println("File write", time.Since(now))
}

func PrepareUncommonCase() *excelize.File {
	var datasource []utils.SalesStatisticalAnalysisItem

	data, _ := os.Open("driver.json")

	byteJson, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteJson, &datasource)

	columns := []utils.ColumnType{
		{
			Field: "BusinessID",
			Title: "사업자명",
		},
		{
			Field: "Name",
			Title: "기사",
		},
		{
			Field: "CallAppType",
			Title: "영업구분",
		},
		{
			Title: "합계",
			Children: []utils.ColumnType{
				{
					Field: "TotalSalesAmount",
					Title: "금액",
					Render: func(v interface{}) interface{} {
						price := v.(string)
						if len(price) < 1 {
							price = "0"
						}

						return price
					},
				},
				{
					Field: "TotalSalesCount",
					Title: "건수",
				},
			},
		},
		{
			Title: "20년 12월",
			Children: []utils.ColumnType{
				{
					Field: "SalesAmount1",
					Title: "금액",
					Render: func(v interface{}) interface{} {
						price := v.(string)
						if len(price) < 1 {
							price = "0"
						}

						return price
					},
				},
				{
					Field: "SalesCount1",
					Title: "건수",
				},
				{
					Field: "WorkDayCount1",
					Title: "근무일수",
					Render: func(v interface{}) interface{} {
						return fmt.Sprintf("%s일", v.(string))
					},
				},
			},
		},
		{
			Title: "21년 01월",
			Children: []utils.ColumnType{
				{
					Field: "SalesAmount2",
					Title: "금액",
					Render: func(v interface{}) interface{} {
						price := v.(string)
						if len(price) < 1 {
							price = "0"
						}

						return price
					},
				},
				{
					Field: "SalesCount2",
					Title: "건수",
				},
				{
					Field: "WorkDayCount2",
					Title: "근무일수",
					Render: func(v interface{}) interface{} {
						return fmt.Sprintf("%s일", v.(string))
					},
				},
			},
		},
	}

	f, _ := utils.
		NewExcel().
		SetDataSource(datasource).
		SetColumns(columns).
		RenderAutoMerging()

	return f
}
