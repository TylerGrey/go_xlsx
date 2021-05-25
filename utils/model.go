package utils

type UserInfo struct {
	ID             string `json:"ID" col_name:"ID"`
	Name           string `json:"name" col_name:"Name" col_order:"4"`
	MobileNumber   string `json:"mobileNumber" col_name:"MobileNumber" col_order:"1"`
	EmployeeRegNum string `json:"employeeRegNum" col_name:"EmployeeRegNum"`
	TeamName       string `json:"teamName" col_name:"TeamName" col_order:"3"`
	CompanyEmail   string `json:"companyEmail" col_name:"CompanyEmail" col_order:"2"`
}

type SalesStatisticalAnalysisItem struct {
	ID               uint   `json:"ID"` // 고유 ID
	BusinessID       string `json:"businessID"`
	TaxiDriverID     string `json:"taxiDriverID"`
	VIN              string `json:"vin"`            // 차대번호
	Name             string `json:"taxiDriverName"` // 기사명
	Pincode          string `json:"pincode"`
	IsCall           string `json:"isCall"`           // 호출여부
	CallAppType      string `json:"callAppType"`      // 호출앱
	TotalSalesAmount string `json:"totalSalesAmount"` // 매출금액
	TotalSalesCount  string `json:"totalSalesCount"`  // 총영업건수
	SalesAmount1     string `json:"salesAmount1"`     // 매출금액
	SalesAmount2     string `json:"salesAmount2"`     // 매출금액
	SalesAmount3     string `json:"salesAmount3"`     // 매출금액
	SalesAmount4     string `json:"salesAmount4"`     // 매출금액
	SalesAmount5     string `json:"salesAmount5"`     // 매출금액
	SalesAmount6     string `json:"salesAmount6"`     // 매출금액
	SalesAmount7     string `json:"salesAmount7"`     // 매출금액
	SalesAmount8     string `json:"salesAmount8"`     // 매출금액
	SalesAmount9     string `json:"salesAmount9"`     // 매출금액
	SalesAmount10    string `json:"salesAmount10"`    // 매출금액
	SalesAmount11    string `json:"salesAmount11"`    // 매출금액
	SalesAmount12    string `json:"salesAmount12"`    // 매출금액
	SalesAmount13    string `json:"salesAmount13"`    // 매출금액
	SalesAmount14    string `json:"salesAmount14"`    // 매출금액
	SalesAmount15    string `json:"salesAmount15"`    // 매출금액
	SalesAmount16    string `json:"salesAmount16"`    // 매출금액
	SalesAmount17    string `json:"salesAmount17"`    // 매출금액
	SalesAmount18    string `json:"salesAmount18"`    // 매출금액
	SalesAmount19    string `json:"salesAmount19"`    // 매출금액
	SalesAmount20    string `json:"salesAmount20"`    // 매출금액
	SalesAmount21    string `json:"salesAmount21"`    // 매출금액
	SalesAmount22    string `json:"salesAmount22"`    // 매출금액
	SalesAmount23    string `json:"salesAmount23"`    // 매출금액
	SalesAmount24    string `json:"salesAmount24"`    // 매출금액
	SalesAmount25    string `json:"salesAmount25"`    // 매출금액
	SalesAmount26    string `json:"salesAmount26"`    // 매출금액
	SalesAmount27    string `json:"salesAmount27"`    // 매출금액
	SalesAmount28    string `json:"salesAmount28"`    // 매출금액
	SalesAmount29    string `json:"salesAmount29"`    // 매출금액
	SalesAmount30    string `json:"salesAmount30"`    // 매출금액
	SalesAmount31    string `json:"salesAmount31"`    // 매출금액
	SalesCount1      string `json:"salesCount1"`      // 영업건수
	SalesCount2      string `json:"salesCount2"`      // 영업건수
	SalesCount3      string `json:"salesCount3"`      // 영업건수
	SalesCount4      string `json:"salesCount4"`      // 영업건수
	SalesCount5      string `json:"salesCount5"`      // 영업건수
	SalesCount6      string `json:"salesCount6"`      // 영업건수
	SalesCount7      string `json:"salesCount7"`      // 영업건수
	SalesCount8      string `json:"salesCount8"`      // 영업건수
	SalesCount9      string `json:"salesCount9"`      // 영업건수
	SalesCount10     string `json:"salesCount10"`     // 영업건수
	SalesCount11     string `json:"salesCount11"`     // 영업건수
	SalesCount12     string `json:"salesCount12"`     // 영업건수
	SalesCount13     string `json:"salesCount13"`     // 영업건수
	SalesCount14     string `json:"salesCount14"`     // 영업건수
	SalesCount15     string `json:"salesCount15"`     // 영업건수
	SalesCount16     string `json:"salesCount16"`     // 영업건수
	SalesCount17     string `json:"salesCount17"`     // 영업건수
	SalesCount18     string `json:"salesCount18"`     // 영업건수
	SalesCount19     string `json:"salesCount19"`     // 영업건수
	SalesCount20     string `json:"salesCount20"`     // 영업건수
	SalesCount21     string `json:"salesCount21"`     // 영업건수
	SalesCount22     string `json:"salesCount22"`     // 영업건수
	SalesCount23     string `json:"salesCount23"`     // 영업건수
	SalesCount24     string `json:"salesCount24"`     // 영업건수
	SalesCount25     string `json:"salesCount25"`     // 영업건수
	SalesCount26     string `json:"salesCount26"`     // 영업건수
	SalesCount27     string `json:"salesCount27"`     // 영업건수
	SalesCount28     string `json:"salesCount28"`     // 영업건수
	SalesCount29     string `json:"salesCount29"`     // 영업건수
	SalesCount30     string `json:"salesCount30"`     // 영업건수
	SalesCount31     string `json:"salesCount31"`     // 영업건수
	WorkDayCount1    string `json:"workDayCount1"`    // 근무일수
	WorkDayCount2    string `json:"workDayCount2"`    // 근무일수
	WorkDayCount3    string `json:"workDayCount3"`    // 근무일수
	WorkDayCount4    string `json:"workDayCount4"`    // 근무일수
	WorkDayCount5    string `json:"workDayCount5"`    // 근무일수
	WorkDayCount6    string `json:"workDayCount6"`    // 근무일수
	WorkDayCount7    string `json:"workDayCount7"`    // 근무일수
	WorkDayCount8    string `json:"workDayCount8"`    // 근무일수
	WorkDayCount9    string `json:"workDayCount9"`    // 근무일수
	WorkDayCount10   string `json:"workDayCount10"`   // 근무일수
	WorkDayCount11   string `json:"workDayCount11"`   // 근무일수
	WorkDayCount12   string `json:"workDayCount12"`   // 근무일수
	WorkDayCount13   string `json:"workDayCount13"`   // 근무일수
	WorkDayCount14   string `json:"workDayCount14"`   // 근무일수
	WorkDayCount15   string `json:"workDayCount15"`   // 근무일수
	WorkDayCount16   string `json:"workDayCount16"`   // 근무일수
	WorkDayCount17   string `json:"workDayCount17"`   // 근무일수
	WorkDayCount18   string `json:"workDayCount18"`   // 근무일수
	WorkDayCount19   string `json:"workDayCount19"`   // 근무일수
	WorkDayCount20   string `json:"workDayCount20"`   // 근무일수
	WorkDayCount21   string `json:"workDayCount21"`   // 근무일수
	WorkDayCount22   string `json:"workDayCount22"`   // 근무일수
	WorkDayCount23   string `json:"workDayCount23"`   // 근무일수
	WorkDayCount24   string `json:"workDayCount24"`   // 근무일수
	WorkDayCount25   string `json:"workDayCount25"`   // 근무일수
	WorkDayCount26   string `json:"workDayCount26"`   // 근무일수
	WorkDayCount27   string `json:"workDayCount27"`   // 근무일수
	WorkDayCount28   string `json:"workDayCount28"`   // 근무일수
	WorkDayCount29   string `json:"workDayCount29"`   // 근무일수
	WorkDayCount30   string `json:"workDayCount30"`   // 근무일수
	WorkDayCount31   string `json:"workDayCount31"`   // 근무일수
}

type TaxiDriverTimeAndAttendanceItem struct {
	ID                 string `json:"ID"`
	BusinessID         string `json:"businessID"`
	No                 string `json:"no"`                 // 레코드 Row number
	RowSpan            string `json:"rowSpan"`            // 해당 의 record수
	LicensePlateNumber string `json:"licensePlateNumber"` // 차량번호
	Name               string `json:"name"`               // 기사 이름
	TaxiDriverID       string `json:"taxiDriverID"`       // 기사 ID
	WorkDay            string `json:"workDay"`            // 근무일수(일)
	WorkTime           string `json:"workTime"`           // 근로시간( 시간 분 )
	AverageWorkTime    string `json:"avgWorkTime"`        // 평균근로시간( 시간 분 )
	Pincode            string `json:"pincode"`            // Pincode
	Key                string `json:"key"`
	LogInOutType       string `json:"logInOutType"` // 로그인 아웃 타입(IN/OUT)
	LogInOutDate       string `json:"logInOutDate"` // 로그인 아웃 시간(format: yy.mm.dd hh:mm:ss)
	WorkDate           string `json:"workDate"`     // 근무일 (format: YYYY.mm.dd)
	Vin                string `json:"vin"`
	NameWithUnMasking  string `json:"nameWithUnMasking"` // masking  적용 하지 않은 이름 ( 엑셀 다운로드 용)
}
