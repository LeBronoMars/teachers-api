package handlers

import (
	"net/http"
	"strconv"

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

//get all class
func (handler ScheduleHandler) Index(c *gin.Context) {
	schedules := []m.Schedule{}		
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	roomParam, roomParamExist := c.GetQuery("room")

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

	//section param exist
	if roomParamExist {
		query = query.Where("section = ?", roomParam)
	} 

	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&schedules)
	c.JSON(http.StatusOK, schedules)	
	return
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
				if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
					result := handler.db.Model(&existingscheduleById).Update(&schedule)
					updatedSchedule := m.Schedule{}
					handler.db.Where("id = ?", schedule.Id).First(&updatedSchedule)
					if result.RowsAffected > 0 {
						c.JSON(http.StatusOK, updatedSchedule)
					} else if result.Error != nil {
						respond(http.StatusOK, result.Error.Error(), c, true)
					} else {
						respond(http.StatusOK, "There are no changes detected.", c , true)
					}
				} else {
					if (c.PostForm("for_deletion") == "true") {
						delete := handler.db.Delete(&existingscheduleById)
						if delete.RowsAffected > 0 {
							deletedSchedule := m.Schedule{}
							if handler.db.Unscoped().Where("id = ?", schedule.Id).First(&deletedSchedule).RowsAffected > 0 {
								c.JSON(http.StatusOK, deletedSchedule)
							}
						} else {
							respond(http.StatusBadRequest, delete.Error.Error(), c, true)
						}
					} else {
						respond(http.StatusBadRequest, "Invalid action.", c, true)
					}
				}
			} else {
				if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
					schedule.CreatedBy = GetCreator(c)
					saveResult := handler.db.Create(&schedule)
					if saveResult.RowsAffected > 0 {
						c.JSON(http.StatusCreated, schedule)
					} else {
						deletedSchedule := m.Schedule{}
						if handler.db.Unscoped().Where("id = ?", schedule.Id).First(&deletedSchedule).RowsAffected > 0 {
							c.JSON(http.StatusOK, deletedSchedule)
						}
					}
				} else {
					if (c.PostForm("for_deletion") == "true") {
						existingSchedule := m.Schedule{}
						if handler.db.Unscoped().Where("id = ?", schedule.Id).First(&existingSchedule).RowsAffected > 0 {
							c.JSON(http.StatusOK, existingSchedule)		
						} else {
							schedule.DeletedAt = GetDeletedAt(c)
							c.JSON(http.StatusOK, schedule)	
						}
					} else {
						respond(http.StatusBadRequest, "Invalid action.", c, true)
					}
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

func (handler ScheduleHandler) Show(c *gin.Context) {
	id := c.Param("id")
	schedule := m.Schedule{}
	scheduleQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&schedule)

	if scheduleQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, schedule)
	} else {
		respond(http.StatusNotFound, "Schedule record not found", c, true)
	}
	return
}


