package models

type Subject struct {
	BaseModel
	SubjectName string `json:"subject_name" form:"subject_name"`
	SubjectCode string `json:"subject_code" form:"subject_code"`
	Description string `json:"desciption" form:"desciption"`
	IsSynced bool `json:"is_synched" form:"is_synched"`
}