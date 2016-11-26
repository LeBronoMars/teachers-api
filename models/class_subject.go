package models

type ClassSubject struct {
	BaseModel
	ClassId string `json:"class_id" form:"class_id" binding:"required"`
	SubjectId string `json:"subject_id" form:"subject_id" binding:"required"`
	TeacherId string `json:"teacher_id" form:"teacher_id" binding:"required"`
}