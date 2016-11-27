package models

type QryClassSchools struct {
	ClassId string `json:"class_id"`
	GradeLevel string `json:"grade_level"`
	Section string `json:"section"`
	SchoolYearFrom string `json:"school_year_from"`
	SchoolYearTo string `json:"school_year_to"`
	SchoolId string `json:"school_id"`
	SchoolName string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
	ContactNo string `json:"contact_no"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedBy string `json:"created_by"`
	Remarks string `json:"remarks"`
	DeletedAt string `json:"deleted_at"`
}