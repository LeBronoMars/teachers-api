package models

type ClassSubject struct {
	BaseModel
	ClassId string `json:"class_id" form:"class_id"`
	SubjectId string `json:"subject_id" form:"subject_id"`
}