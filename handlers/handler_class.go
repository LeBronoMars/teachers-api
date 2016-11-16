package handlers

import (
	"net/http"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type ClassHandler struct {
	db *gorm.DB
}

func NewClassHandler(db *gorm.DB) *ClassHandler {
	return &ClassHandler{db}
}

func (handler ClassHandler) Index(c *gin.Context) {
	classess := []m.QryClassSchools{}		
	
	var query = handler.db

	startParam,startParamExist := c.GetQuery("start")
	limitParam,limitParamExist := c.GetQuery("limit")
	orderParam,orderParamExist := c.GetQuery("order")


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

	query.Find(&classess)
	c.JSON(http.StatusOK, classess)	
	return
}

//create new class
func (handler ClassHandler) Create(c *gin.Context) {
	var newClass m.Class
	err := c.Bind(&newClass)
	if err == nil {
		//check if school is existing
		school := m.School{}
		existingSchoolQuery := handler.db.Where("id = ?", newClass.School).First(&school)
		if existingSchoolQuery.RowsAffected > 0 {
			//check if class is already existing
			existingClass := m.Class{}
			existingClassQuery := handler.db.Where("section = ? and grade_level = ?", newClass.Section, newClass.GradeLevel).First(&existingClass)

			if existingClassQuery.RowsAffected == 0 {
				saveResult := handler.db.Save(&newClass)
				if (saveResult.RowsAffected > 0) {
					qryNewClass := m.QryClassSchools{}
					handler.db.Where("class_id = ?", newClass.Id).First(&qryNewClass)
					c.JSON(http.StatusCreated, qryNewClass)
				} else {
					respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
				}
			} else {
				respond(http.StatusBadRequest, fmt.Sprintf("Class with section of %s in Grade Level %v already exist.", newClass.Section, newClass.GradeLevel), c, true)
			}
		} else {
			respond(http.StatusBadRequest, "School not found.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}