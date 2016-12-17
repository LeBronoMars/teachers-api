package models

type Attendance struct {
	BaseModel
	ClassStudentId string `json:"class_student_id" form:"class_student_id" binding:"required"`
	Attendance string `json:"attendance" form:"attendance" binding:"required"`
	ScheduleId string `json:"schedule_id" form:"schedule_id" binding:"required"`
} 
