package models

import (
	"crypto/md5"
	"fmt"
)

type User struct {
	BaseModel
	EmployeeNo string `json:"employee_no" form:"employee_no"`
	FirstName  string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`	
	MiddleName string `json:"middle_name" form:"middle_name"`
	BirthDate string `json:"birth_date" form:"birth_date"`
	BirthPlace string `json:"birth_place" form:"birth_place"`
	Gender string `json:"gender" form:"gender"`
	CivilStatus string `json:"civil_status" form:"civil_status"`
	Email string `json:"email" form:"email"`
	Address string `json:"address" form:"address"`
	ContactNo string `json:"contact_no" form:"contact_no"`
	Status string `json:"status"`
	UserRole string `json:"user_role" form:"user_role"`
	Password string `json:"-" form:"password"`
	PicUrl string `json:"pic_url" form:"pic_url"`
	Position string `json:"position" form:"position"`
	SchoolId string `json:"school_id" form:"school_id"`
}

func (u *User) BeforeCreate() (err error) {
	u.Status = "active"
	u.IsSynced = true
	defaultPic := fmt.Sprintf("%x", md5.Sum([]byte(u.Email)))
	u.PicUrl = fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon", defaultPic)
	return
}

func (u *User) AfterUpdate() (err error) {
	u.IsSynced = true
	return
}