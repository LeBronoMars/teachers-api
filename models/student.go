package models

import (
	"crypto/md5"
	"fmt"
)

type Student struct {
	BaseModel
	StudentNo string `json:"student_no" form:"student_no"`
	FirstName string `json:"first_name" form:"first_name"`
	MiddleName string `json:"middle_name" form:"middle_name"`
	LastName string `json:"last_name" form:"last_name"`
	BirthDate string `json:"birth_date" form:"birth_date"`
	Gender string `json:"gender" form:"gender"`
	Address string `json:"address" form:"address"`
	PicUrl string `json:"pic_url" form:"pic_url"`
	Status string `json:"status" form:"status"`
	Remarks string `json:"remarks" form:"remarks"`
}

func (s *Student) BeforeCreate() (err error) {
	s.Status = "active"
	s.IsSynced = true
	defaultPic := fmt.Sprintf("%x", md5.Sum([]byte(s.StudentNo)))
	s.PicUrl = fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon", defaultPic)
	return
}

func (s *Student) BeforeUpdate() (err error) {
	s.IsSynced = true
	return
}

