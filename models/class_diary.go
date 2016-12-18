package models

type ClassDiary struct {
	BaseModel
	Date string `json:"date" form:"date" binding:"required"`
	EventActivity string `json:"event_activity" form:"event_activity" binding:"required"`
	Comment string `json:"comment" form:"comment" binding:"required"`
	ClassId string `json:"class_id" form:"class_id" binding:"required"`
} 