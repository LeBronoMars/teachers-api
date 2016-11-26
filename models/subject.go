package models

type Subject struct {
	BaseModel
	SubjectName string `json:"subject_name" form:"subject_name" binding:"required"`
	SubjectCode string `json:"subject_code" form:"subject_code" binding:"required"`
	Description string `json:"desciption" form:"desciption"`
}