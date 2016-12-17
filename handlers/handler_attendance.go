package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type AttendanceHandler struct {
	db *gorm.DB
}

func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{db}
}

//get all schedule
func (handler AttendanceHandler) Index(c *gin.Context) {
	attendance := []m.Attendance{}		
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	studentIdParam, studentIdParamExist := c.GetQuery("student_id")

	//start param exist
	if startParamExist {
		start,_ := strconv.Atoi(startParam)
		if start != 0 {
			query = query.Offset(start)				
		} else {
			query = query.Offset(0)
		}
	} 

	//limit param exist
	if limitParamExist {
		limit,_ := strconv.Atoi(limitParam)
		query = query.Limit(limit)
	} else {
		query = query.Limit(10)
	}

	//sort param exist
	if orderParamExist {
		query = query.Order(orderParam)
	} 

	if studentIdParamExist {
		query = query.Where("class_student_id = ?", studentIdParam)
	} 

	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&attendance)
	c.JSON(http.StatusOK, attendance)	
	return
}

//create new attendance
func (handler AttendanceHandler) Create(c *gin.Context) {
	var attendance m.Attendance
	err := c.Bind(&attendance)

	if err == nil {
		//check if class student is existing
		existingClassStudent := m.ClassStudent{}
		existingClassStudentQuery := handler.db.Where("id = ?", attendance.ClassStudentId).First(&existingClassStudent)

		if existingClassStudentQuery.RowsAffected > 0 {

			//check if schedule exists
			existingSchedule := m.Schedule{}
			if handler.db.Where("id = ?", attendance.ScheduleId).First(&existingSchedule).RowsAffected > 0 {
				existingAttendanceById := m.Attendance{}
				if handler.db.Where("id = ?", attendance.Id).First(&existingAttendanceById).RowsAffected > 0 {
					result := handler.db.Model(&existingAttendanceById).Update(&attendance)
					if result.RowsAffected > 0 {
						c.JSON(http.StatusOK, attendance)
					} else if result.Error != nil {
						respond(http.StatusBadRequest, result.Error.Error(), c, true)
					} else {
						respond(http.StatusBadRequest, "There are no changes detected.", c , true)
					}
				} else {
					attendance.CreatedBy = GetCreator(c)
					saveResult := handler.db.Create(&attendance)
					if saveResult.RowsAffected > 0 {
						c.JSON(http.StatusCreated, attendance)
					} else {
						respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
					}
				}
			} else {
				respond(http.StatusBadRequest, "Schedule record not found", c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Class student record not found", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler AttendanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	attendance := m.Attendance{}
	attendanceQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&attendance)

	if attendanceQuery.RowsAffected > 0 {
		deleteResult := handler.db.Delete(&attendance)
		if deleteResult.RowsAffected > 0 {
			respond(http.StatusOK, "attendance record successfully deleted", c, false)
		} else {
			respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
		}
	} else {
		respond(http.StatusNotFound, "attendance record not found", c, true)
	}
	return
}

func (handler AttendanceHandler) Show(c *gin.Context) {
	id := c.Param("id")
	attendance := m.Attendance{}
	attendanceQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&attendance)

	if attendanceQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, attendance)
	} else {
		respond(http.StatusNotFound, "attendance record not found", c, true)
	}
	return
}

