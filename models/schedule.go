package models

type Schedule struct {
	BaseModel
	Date string `json:"date" form:"date" binding:"required"`
	StartTime string `json:"start_time" form:"start_time" binding:"required"`
	EndTime string `json:"end_time" form:"end_time" binding:"required"`
	ClassRoom string `json:"class_room" form:"class_room" binding:"required"`
	ClassSubjectId string `json:"class_subject_id" form:"class_subject_id" binding:"required"`
}
