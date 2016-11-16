package models

import "time"

type Class struct {
	BaseModel
	GradeLevel string `json:"grade_level" form:"grade_level" binding:"required"`
	Section string `json:"section" form:"section" binding:"required"`
	SchoolYear int `json:"school_year"`
	School int `json:"school_id" form:"school_id" binding:"required"`
	Remarks string `json:"remarks" form:"remarks"`
	IsSynced bool `json:"is_synced" form:"is_synced"`
}

func (c *Class) BeforeCreate() (err error) {
	c.SchoolYear = time.Now().Year()
	c.IsSynced = true
	return
}