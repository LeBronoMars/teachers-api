package models

import "errors"

type Class struct {
	BaseModel
	GradeLevel string `json:"grade_level" form:"grade_level" binding:"required"`
	Section string `json:"section" form:"section" binding:"required"`
	SchoolYearFrom int `json:"school_year_from" form:"school_year_from" binding:"required"`
	SchoolYearTo int `json:"school_year_to" form:"school_year_to" binding:"required"`
	School string `json:"school_id" form:"school_id" binding:"required"`
	Remarks string `json:"remarks" form:"remarks"`
}

func (c *Class) BeforeCreate() (err error) {
	c.IsSynced = true
	if c.SchoolYearFrom >= c.SchoolYearTo {
		err = errors.New("School year to must be greater than school year from")
	} else if ((c.SchoolYearTo - c.SchoolYearFrom) > 1) {
		err = errors.New("Invalid school year duration.")
	}
	return
}