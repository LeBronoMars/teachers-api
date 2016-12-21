package models

type ClassStudent struct {
	BaseModel
	Class Class `json:"class"`
	ClassId string `json:"-" form:"class_id" sql:"type:varchar(100)"`
	Student Student `json:"student"`
	StudentId string `json:"-" form:"student_id" sql:"type:varchar(100)"`
}