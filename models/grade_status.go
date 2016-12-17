package models

type GradeStatus struct {
	BaseModel
	ClassSubjectId string `json:"class_subject_id" form:"class_subject_id" binding:"required"`
	StudentId string `json:"student_id" form:"student_id" binding:"required"`
	IsPassed string `json:"is_passed" form:"is_passed"`
} 
