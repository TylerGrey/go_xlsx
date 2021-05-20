package utils

type TestType struct {
	A string `col_name:"Col A"`
	B string `col_name:"-"`
	C string `col_name:"Col B"`
	D string
	E string `col_name:"Col E"`
	F string `col_name:"Col F"`
	G string `col_name:"Col G"`
	H string `col_name:"Col H"`
	I string `col_name:"Col I"`
	J string `col_name:"Col J"`
}

type UserInfo struct {
	ID             string `json:"ID" col_name:"ID"`
	Name           string `json:"name" col_name:"Name" col_order:"4"`
	MobileNumber   string `json:"mobileNumber" col_name:"MobileNumber" col_order:"1"`
	EmployeeRegNum string `json:"employeeRegNum" col_name:"EmployeeRegNum"`
	TeamName       string `json:"teamName" col_name:"TeamName" col_order:"3"`
	CompanyEmail   string `json:"companyEmail" col_name:"CompanyEmail" col_order:"2"`
}
