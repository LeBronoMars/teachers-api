package models

type SchoolYear struct {
	BaseModel
	From  int `json:"from" form:"from"`
	To  int `json:"to" form:"to"`
	Description  string `json:"description" form:"description"`	
}