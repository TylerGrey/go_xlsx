package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"xlsx_test/utils"
)

func main() {
	// File
	now := time.Now()
	PrepareAndReturnExcel().SaveAs("Test.xlsx")
	fmt.Println(time.Since(now))

	// HTTP
	//http.HandleFunc("/xlsx", downloadExcel)
	//http.ListenAndServe(":3000", nil)
}

func PrepareAndReturnExcel() *excelize.File {
	var datasource []utils.UserInfo
	for i := 0; i < 1_000_000; i++ {
		datasource = append(datasource, utils.UserInfo{
			ID:             strconv.Itoa(i),
			Name:           fmt.Sprintf("Name%d", i),
			MobileNumber:   fmt.Sprintf("010-%4d-%4d", rand.Intn(9999), rand.Intn(9999)),
			EmployeeRegNum: strconv.Itoa(i),
			TeamName:       fmt.Sprintf("Team %d", rand.Int()),
			CompanyEmail:   fmt.Sprintf("%d@gmail.com", rand.Int()),
		})
		//datasource = append(datasource, utils.TestType{
		//	A: "A1",
		//	B: "B1",
		//	C: "C1",
		//	D: "D1",
		//	E: "E1",
		//	F: "F1",
		//	G: "G1",
		//	H: "H1",
		//	I: "I1",
		//	J: "J1",
		//})
	}

	f, _ := utils.
		NewExcel().
		SetSheet("Test Sheet").
		SetStartRow(2).
		SetDataSource(datasource).
		Render()

	return f
}

func downloadExcel(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	// Get the Excel file with the user input data
	file := PrepareAndReturnExcel()

	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=userInputData.xlsx")
	w.Header().Set("File-Name", "userInputData.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	file.Write(w)
	fmt.Println(time.Since(now))
}
