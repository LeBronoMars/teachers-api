package models

import (
	"crypto/md5"
	"fmt"
)

type Teacher struct {
	BaseModel
	EmployeeNo string `json:"employee_no" form:"employee_no" binding:"required"`
	FirstName  string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`	
	MiddleName string `json:"middle_name" form:"middle_name" binding:"required"`
	BirthDate string `json:"birth_date" form:"birth_date" binding:"required"`
	Gender string `json:"gender" form:"gender" binding:"required"`
	Email string `json:"email" form:"email" binding:"required"`
	Address string `json:"address" form:"address" binding:"required"`
	ContactNo string `json:"contact_no" form:"contact_no" binding:"required"`
	Status string `json:"status"`
	PicUrl string `json:"pic_url" form:"pic_url"`
	Position string `json:"position" form:"position" binding:"required"`
	IsSynced bool `json:"is_synced"`
	School int `json:"school"`
}

func (t *Teacher) BeforeCreate() (err error) {
	t.Status = "active"
	t.IsSynced = true
	defaultPic := fmt.Sprintf("%x", md5.Sum([]byte(t.Email)))
	t.PicUrl = fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon", defaultPic)
	return
}