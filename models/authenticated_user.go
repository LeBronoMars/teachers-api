package models

import "time"

type AuthenticatedUser struct {
	Id int `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName  string `json:"middle_name"`	
	LastName  string `json:"last_name"`	
	Status string `json:"status"`
	UserLevel string `json:"user_level"`
	Email string `json:"email"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Gender string `json:"gender"`
	PicUrl string `json:"pic_url"`
	Token string `json:"token"`
	Position string `json:"position"`
	ContactNo string `json:"contact_no"`
}