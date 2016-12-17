package models

type Schedule struct {
	BaseModel
	Date string `json:"date" form:"date"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime string `json:"end_time" form:"end_time"`
	ClassRoom string `json:"class_room" form:"class_room"`
	ClassSubjectId string `json:"class_subject_id" form:"class_subject_id"`
}
