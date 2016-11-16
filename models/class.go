package models

import "time"

type Class struct {
	BaseModel
	GradeLevel int `json:"grade_level" form:"grade_level" binding:"required"`
	Section string `json:"section" form:"section" binding:"required"`
	SchoolYear int `json:"school_year"`
	SchoolId int `json:"school_id" form:"school_id" binding:"required"`
	isSynced bool `json:"is_synced" form:"is_synced"`
}

func (c *Class) BeforeCreate() (err error) {
	c.SchoolYear = time.Now().Year()
	c.isSynced = true
	return
}