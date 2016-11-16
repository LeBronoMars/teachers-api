package models

type QryClassSchools struct {
	ClassId int `json:"class_id"`
	GradeLevel string `json:"grade_level"`
	Section string `json:"section"`
	SchoolYear string `json:"school_year"`
	SchoolId int `json:"school_id"`
	SchoolName string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
	ContactNo string `json:"contact_no"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

}