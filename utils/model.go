package utils

type TestType struct {
	A string `col_name:"Col A"`
	B string `col_name:"-"`
	C string `col_name:"Col B"`
	D string
}

type UserInfo struct {
	ID             string `json:"ID" col_name:"ID" col_order:"1"`
	Name           string `json:"name" col_name:"Name" col_order:"2"`
	MobileNumber   string `json:"mobileNumber" col_name:"MobileNumber" col_order:"3"`
	EmployeeRegNum string `json:"employeeRegNum" col_name:"EmployeeRegNum" col_order:"4"`
	TeamName       string `json:"teamName" col_name:"TeamName" col_order:"5"`
	CompanyEmail   string `json:"companyEmail" col_name:"CompanyEmail" col_order:"6"`
}