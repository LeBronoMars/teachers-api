package models

import "time"

type Class struct {
	BaseModel
	GradeLevel string `json:"grade_level" form:"grade_level" binding:"required"`
	Section string `json:"section" form:"section" binding:"required"`
	SchoolYearFrom int `json:"school_year_from"`
	SchoolYearTo int `json:"school_year_to"`
	School string `json:"school_id" form:"school_id" binding:"required"`
	Remarks string `json:"remarks" form:"remarks"`
}

func (c *Class) BeforeCreate() (err error) {
	c.SchoolYearFrom = time.Now().Year()
	c.SchoolYearTo = time.Now().Year() + 1
	c.IsSynced = true
	return
}