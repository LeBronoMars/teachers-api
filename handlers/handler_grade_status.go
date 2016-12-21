package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type GradeStatusHandler struct {
	db *gorm.DB
}

func NewGradeStatusHandler(db *gorm.DB) *GradeStatusHandler {
	return &GradeStatusHandler{db}
}

//get all schedule
func (handler GradeStatusHandler) Index(c *gin.Context) {
	gradeStatus := []m.GradeStatus{}		
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	classSubjectParam, classSubjectParamExist := c.GetQuery("class_subject_id")
	statusParam, statusParamExist := c.GetQuery("is_passed")

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

	if classSubjectParamExist {
		query = query.Where("class_subject_id = ?", classSubjectParam)
	} 

	if statusParamExist {
		query = query.Where("is_passed = ?", statusParam)
	} 

	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&gradeStatus)
	c.JSON(http.StatusOK, gradeStatus)	
	return
}

//create new attendance
func (handler GradeStatusHandler) Create(c *gin.Context) {
	var gradeStatus m.GradeStatus
	err := c.Bind(&gradeStatus)

	if err == nil {
		//check if class subject is existing
		existingClassSubject := m.ClassSubject{}
		existingClassSubjectQuery := handler.db.Where("id = ?", gradeStatus.ClassSubjectId).First(&existingClassSubject)

		if existingClassSubjectQuery.RowsAffected > 0 {

			//check if student exists
			existingStudent := m.Student{}
			if handler.db.Where("id = ?", gradeStatus.StudentId).First(&existingStudent).RowsAffected > 0 {
				existingGradeStatus := m.GradeStatus{}
				if handler.db.Where("id = ?", gradeStatus.Id).First(&existingGradeStatus).RowsAffected > 0 {
					if (c.PostForm("for_deletion") == "") {
						result := handler.db.Model(&existingGradeStatus).Update(&gradeStatus)
						if result.RowsAffected > 0 {
							c.JSON(http.StatusOK, gradeStatus)
						} else if result.Error != nil {
							respond(http.StatusBadRequest, result.Error.Error(), c, true)
						} else {
							respond(http.StatusBadRequest, "There are no changes detected.", c , true)
						}
					} else {
						delete := handler.db.Delete(&existingGradeStatus)
						if delete.RowsAffected > 0 {
							respond(http.StatusOK, "Record successfully deleted.", c, false)
						} else {
							respond(http.StatusBadRequest, delete.Error.Error(), c, true)
						}
					}
				} else {
					gradeStatus.CreatedBy = GetCreator(c)
					saveResult := handler.db.Create(&gradeStatus)
					if saveResult.RowsAffected > 0 {
						c.JSON(http.StatusCreated, gradeStatus)
					} else {
						respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
					}
				}
			} else {
				respond(http.StatusBadRequest, "Student record not found", c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Class subject record not found", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler GradeStatusHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	gradeStatus := m.GradeStatus{}
	gradeStatusQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&gradeStatus)

	if gradeStatusQuery.RowsAffected > 0 {
		deleteResult := handler.db.Delete(&gradeStatus)
		if deleteResult.RowsAffected > 0 {
			respond(http.StatusOK, "Grade status record successfully deleted", c, false)
		} else {
			respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
		}
	} else {
		respond(http.StatusNotFound, "Grade status record not found", c, true)
	}
	return
}

func (handler GradeStatusHandler) Show(c *gin.Context) {
	id := c.Param("id")
	gradeStatus := m.GradeStatus{}
	gradeStatusQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&gradeStatus)

	if gradeStatusQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, gradeStatus)
	} else {
		respond(http.StatusNotFound, "Grade status record not found", c, true)
	}
	return
}


