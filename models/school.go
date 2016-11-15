package models

type School struct {
	BaseModel
	SchoolName  string `json:"school_name" form:"school_name" binding:"required"`
	SchoolAddress  string `json:"school_address" form:"school_address" binding:"required"`	
	ContactNo string `json:"contact_no" form:"contact_no" binding:"required"`
	
}