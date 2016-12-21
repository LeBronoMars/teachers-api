package models

type ClassStudent struct {
	BaseModel
	ClassId string `json:"class_id" form:"class_id" sql:"type:varchar(100)"`
	StudentId string `json:"student_id" form:"student_id" sql:"type:varchar(100)"`
}