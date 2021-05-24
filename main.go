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
	f := PrepareAndReturnExcel()

	now := time.Now()
	if err := f.SaveAs("Test.xlsx"); err != nil {
		return
	}
	fmt.Println("SaveAs", time.Since(now))

	// HTTP
	//http.HandleFunc("/xlsx", downloadExcel)
	//http.ListenAndServe(":3000", nil)

	// Uncommon case
	//PrepareUncommonCase().SaveAs("merge.xlsx")
}

func PrepareAndReturnExcel() *excelize.File {
	now := time.Now()
	var datasource []*utils.UserInfo
	for i := 0; i < 1_000; i++ {
		datasource = append(datasource, &utils.UserInfo{
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
		SetAutoMerge(true).
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
	var datasource []*utils.SalesStatisticalAnalysisItem

	//driverData, _ := os.Open("driver.json")
	//byteJson, _ := ioutil.ReadAll(driverData)
	//json.Unmarshal(byteJson, &datasource)
	//
	//columns := []utils.ColumnType{
	//	{
	//		Field: "BusinessID",
	//		Name: "사업자명",
	//		MergeColumn: true,
	//	},
	//	{
	//		Field: "Name",
	//		Name: "기사",
	//		MergeColumn: true,
	//	},
	//	{
	//		Field: "CallAppType",
	//		Name: "영업구분",
	//	},
	//	{
	//		Name: "합계",
	//		Children: []utils.ColumnType{
	//			{
	//				Field: "TotalSalesAmount",
	//				Name: "금액",
	//				Render: func(v interface{}) interface{} {
	//					price := v.(string)
	//					if len(price) < 1 {
	//						price = "0"
	//					}
	//
	//					return price
	//				},
	//			},
	//			{
	//				Field: "TotalSalesCount",
	//				Name: "건수",
	//			},
	//		},
	//	},
	//	{
	//		Name: "20년 12월",
	//		Children: []utils.ColumnType{
	//			{
	//				Field: "SalesAmount1",
	//				Name: "금액",
	//				Render: func(v interface{}) interface{} {
	//					price := v.(string)
	//					if len(price) < 1 {
	//						price = "0"
	//					}
	//
	//					return price
	//				},
	//			},
	//			{
	//				Field: "SalesCount1",
	//				Name: "건수",
	//			},
	//			{
	//				Field: "WorkDayCount1",
	//				Name: "근무일수",
	//				MergeColumn: true,
	//				Render: func(v interface{}) interface{} {
	//					return fmt.Sprintf("%s일", v.(string))
	//				},
	//			},
	//		},
	//	},
	//	{
	//		Name: "21년 01월",
	//		Children: []utils.ColumnType{
	//			{
	//				Field: "SalesAmount2",
	//				Name: "금액",
	//				Render: func(v interface{}) interface{} {
	//					price := v.(string)
	//					if len(price) < 1 {
	//						price = "0"
	//					}
	//
	//					return price
	//				},
	//			},
	//			{
	//				Field: "SalesCount2",
	//				Name: "건수",
	//			},
	//			{
	//				Field: "WorkDayCount2",
	//				Name: "근무일수",
	//				MergeColumn: true,
	//				Render: func(v interface{}) interface{} {
	//					return fmt.Sprintf("%s일", v.(string))
	//				},
	//			},
	//		},
	//	},
	//}

	taxiData, _ := os.Open("taxi.json")
	byteJson, _ := ioutil.ReadAll(taxiData)
	json.Unmarshal(byteJson, &datasource)

	columns := []utils.ColumnType{
		{
			Field:       "BusinessID",
			Name:        "사업자명",
			MergeColumn: true,
		},
		{
			Field: "CallAppType",
			Name:  "영업구분",
		},
		{
			Name: "합계",
			Children: []utils.ColumnType{
				{
					Field: "TotalSalesAmount",
					Name:  "금액",
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
					Name:  "건수",
				},
			},
		},
		{
			Name: "20년 12월",
			Children: []utils.ColumnType{
				{
					Field: "SalesAmount1",
					Name:  "금액",
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
					Name:  "건수",
				},
			},
		},
		{
			Name: "21년 01월",
			Children: []utils.ColumnType{
				{
					Field: "SalesAmount2",
					Name:  "금액",
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
					Name:  "건수",
				},
			},
		},
	}

	f, _ := utils.
		NewExcel().
		SetStartRow(2).
		SetDataSource(datasource).
		SetColumns(columns).
		SetAutoMerge(true).
		Render()

	return f
}
