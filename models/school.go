package models

type School struct {
	BaseModel
	SchoolName  string `json:"school_name" form:"school_name" binding:"required"`
	SchoolAddress  string `json:"school_address" form:"school_address" binding:"required"`	
	ContactNo string `json:"contact_no" form:"contact_no" binding:"required"`
	Latitude float64 `json:"latitude" form:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" form:"longitude" binding:"required"`
}

func (s *School) BeforeCreate() (err error) {
	s.IsSynced = true
	return
}