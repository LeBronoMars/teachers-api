package models

import (
	"crypto/md5"
	"fmt"
)

type User struct {
	BaseModel
	EmployeeNo string `json:"employee_no" form:"employee_no" binding:"required"`
	FirstName  string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`	
	MiddleName string `json:"middle_name" form:"middle_name" binding:"required"`
	BirthDate string `json:"birth_date" form:"birth_date" binding:"required"`
	BirthPlace string `json:"birth_place" form:"birth_place" binding:"required"`
	Gender string `json:"gender" form:"gender" binding:"required"`
	CivilStatus string `json:"civil_status" form:"civil_status"`
	Email string `json:"email" form:"email" binding:"required"`
	Address string `json:"address" form:"address" binding:"required"`
	ContactNo string `json:"contact_no" form:"contact_no" binding:"required"`
	Status string `json:"status"`
	UserRole string `json:"user_role" form:"user_role" binding:"required"`
	Password string `json:"-" form:"password" binding:"required"`
	PicUrl string `json:"pic_url" form:"pic_url"`
	Position string `json:"position" form:"position" binding:"required"`
	SchoolId string `json:"school_id" form:"school_id" binding:"required"`
}

func (u *User) BeforeCreate() (err error) {
	u.Status = "active"
	u.IsSynced = true
	defaultPic := fmt.Sprintf("%x", md5.Sum([]byte(u.Email)))
	u.PicUrl = fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon", defaultPic)
	return
}