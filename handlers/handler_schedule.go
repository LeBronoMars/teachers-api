package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type ScheduleHandler struct {
	db *gorm.DB
}

func NewScheduleHandler(db *gorm.DB) *ScheduleHandler {
	return &ScheduleHandler{db}
}

//create new schedule
func (handler ScheduleHandler) Create(c *gin.Context) {
	var schedule m.Schedule
	err := c.Bind(&schedule)

	if err == nil {
		//check if class subject is existing
		existingClassSubject := m.ClassSubject{}
		existingClassSubjectQuery := handler.db.Where("id = ?", schedule.ClassSubjectId).First(&existingClassSubject)

		if existingClassSubjectQuery.RowsAffected > 0 {
			existingscheduleById := m.Schedule{}
			if handler.db.Where("id = ?", schedule.Id).First(&existingscheduleById).RowsAffected > 0 {
				result := handler.db.Model(&existingscheduleById).Update(&schedule)
				if result.RowsAffected > 0 {
					c.JSON(http.StatusOK, schedule)
				} else if result.Error != nil {
					respond(http.StatusBadRequest, result.Error.Error(), c, true)
				} else {
					respond(http.StatusBadRequest, "There are no changes detected.", c , true)
				}
			} else {
				schedule.CreatedBy = GetCreator(c)
				saveResult := handler.db.Create(&schedule)
				if saveResult.RowsAffected > 0 {
					c.JSON(http.StatusCreated, schedule)
				} else {
					respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
				}
			}
		} else {
			respond(http.StatusBadRequest, "Class subject record not found", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler ScheduleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	schedule := m.Schedule{}
	scheduleQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&schedule)

	if scheduleQuery.RowsAffected > 0 {
		deleteResult := handler.db.Delete(&schedule)
		if deleteResult.RowsAffected > 0 {
			respond(http.StatusOK, "Schedule record successfully deleted", c, false)
		} else {
			respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
		}
	} else {
		respond(http.StatusNotFound, "Schedule record not found", c, true)
	}
	return
}




