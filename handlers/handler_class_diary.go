package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type ClassDiaryHandler struct {
	db *gorm.DB
}

func NewClassDiaryHandler(db *gorm.DB) *ClassDiaryHandler {
	return &ClassDiaryHandler{db}
}

//get all schedule
func (handler ClassDiaryHandler) Index(c *gin.Context) {
	classDiary := []m.ClassDiary{}		
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	classId, classIdExist := c.GetQuery("class_id")

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

	if classIdExist {
		query = query.Where("class_id = ?", classId)
	} 

	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&classDiary)
	c.JSON(http.StatusOK, classDiary)	
	return
}

//create new class diary
func (handler ClassDiaryHandler) Create(c *gin.Context) {
	var classDiary m.ClassDiary
	err := c.Bind(&classDiary)

	if err == nil {
		//check if class 
		existingClass := m.Class{}
		existingClassQuery := handler.db.Where("id = ?", classDiary.ClassId).First(&existingClass)

		if existingClassQuery.RowsAffected > 0 {

			//check if class diary exists
			existingclassDiaryById := m.ClassDiary{}
			if handler.db.Where("id = ?", classDiary.Id).First(&existingclassDiaryById).RowsAffected > 0 {
				if (c.PostForm("for_deletion") == "false") {
					result := handler.db.Model(&existingclassDiaryById).Update(&classDiary)
					if result.RowsAffected > 0 {
						updatedClassDiary := m.ClassDiary{}
						handler.db.Where("id = ?", classDiary.Id).First(&updatedClassDiary)
						c.JSON(http.StatusOK, updatedClassDiary)
					} else if result.Error != nil {
						respond(http.StatusBadRequest, result.Error.Error(), c, true)
					} else {
						respond(http.StatusBadRequest, "There are no changes detected.", c , true)
					}
				} else {
					if (c.PostForm("for_deletion") == "true") {
						delete := handler.db.Delete(&existingclassDiaryById)
						if delete.RowsAffected > 0 {
							respond(http.StatusOK, "Record successfully deleted.", c, false)
						} else {
							respond(http.StatusBadRequest, delete.Error.Error(), c, true)
						}
					} else {
						respond(http.StatusBadRequest, "Invalid action.", c, true)
					}
				}
			} else {
				classDiary.CreatedBy = GetCreator(c)
				saveResult := handler.db.Create(&classDiary)
				if saveResult.RowsAffected > 0 {
					c.JSON(http.StatusCreated, classDiary)
				} else {
					respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
				}
			}
		} else {
			respond(http.StatusBadRequest, "Class record not found", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler ClassDiaryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	classDiary := m.ClassDiary{}
	classDiaryQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&classDiary)

	if classDiaryQuery.RowsAffected > 0 {
		deleteResult := handler.db.Delete(&classDiary)
		if deleteResult.RowsAffected > 0 {
			respond(http.StatusOK, "Class diary record successfully deleted", c, false)
		} else {
			respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
		}
	} else {
		respond(http.StatusNotFound, "Class diary record not found", c, true)
	}
	return
}

func (handler ClassDiaryHandler) Show(c *gin.Context) {
	id := c.Param("id")
	classDiary := m.ClassDiary{}
	classDiaryQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", id, GetCreator(c)).First(&classDiary)

	if classDiaryQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, classDiary)
	} else {
		respond(http.StatusNotFound, "Class diary record not found", c, true)
	}
	return
}


