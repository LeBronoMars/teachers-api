package models

type QryClassSchools struct {
	ClassId int `json:"class_id"`
	GradeLevel string `json:"grade_level"`
	Section string `json:"section"`
	SchoolYearFrom string `json:"school_year_from"`
	SchoolYearTo string `json:"school_year_to"`
	SchoolId int `json:"school_id"`
	SchoolName string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
	ContactNo string `json:"contact_no"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Remarks string `json:"remarks"`
}