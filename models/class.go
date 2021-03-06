package models

import "errors"

type Class struct {
	BaseModel
	GradeLevel string `json:"grade_level" form:"grade_level"`
	Section string `json:"section" form:"section"`
	SchoolYearFrom int `json:"school_year_from" form:"school_year_from"`
	SchoolYearTo int `json:"school_year_to" form:"school_year_to"`
	School string `json:"school_id" form:"school_id"`
	Remarks string `json:"remarks" form:"remarks"`
}

func (c *Class) BeforeCreate() (err error) {
	c.IsSynced = true
	if (c.SchoolYearFrom < 2000) {
		err = errors.New("School year from must be year 2000 onwards")
	} else if (c.SchoolYearTo < 2000) {
		err = errors.New("School year to must be year 2000 onwards")
	} else if c.SchoolYearFrom >= c.SchoolYearTo {
		err = errors.New("School year to must be greater than school year from")
	} else if ((c.SchoolYearTo - c.SchoolYearFrom) > 1) {
		err = errors.New("Invalid school year duration.")
	}
	return
}

func (c *Class) BeforeUpdate() (err error) {
	c.IsSynced = true
	return
}