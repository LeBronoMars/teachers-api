package models

import (
	"crypto/md5"
	"fmt"
)

type Student struct {
	BaseModel
	StudentNo string `json:"student_no" form:"student_no" binding:"required"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	MiddleName string `json:"middle_name" form:"middle_name" binding:"required"`
	LastName string `json:"last_name" form:"last_name" binding:"required"`
	BirthDate string `json:"birth_date" form:"birth_date" binding:"required"`
	Gender string `json:"gender" form:"gender" binding:"required"`
	Address string `json:"address" form:"address" binding:"required"`
	PicUrl string `json:"pic_url" form:"pic_url"`
	Status string `json:"status" form:"status"`
	Remarks string `json:"remarks" form:"remarks"`
}

func (s *Student) BeforeCreate() (err error) {
	s.Status = "active"
	defaultPic := fmt.Sprintf("%x", md5.Sum([]byte(s.StudentNo)))
	s.PicUrl = fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon", defaultPic)
	return
}

