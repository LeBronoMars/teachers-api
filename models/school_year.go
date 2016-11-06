package models

type SchoolYear struct {
	BaseModel
	Year  int `json:"year" form:"year"`
	Description  string `json:"description" form:"description"`	
}