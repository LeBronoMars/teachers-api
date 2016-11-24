package models

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
	IsSynced bool `json:"is_synched" form:"is_synched"`
	Status string `json:"status" form:"status"`
	Remarks string `json:"remarks" form:"remarks"`
}

func (s *Student) BeforeCreate() (err error) {
	s.Status = "active"
	return
}

