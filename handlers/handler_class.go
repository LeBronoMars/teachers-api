package handlers

import (
	"net/http"
	"fmt"

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

//create new class
func (handler ClassHandler) Create(c *gin.Context) {
	var newClass m.Class
	err := c.Bind(&newClass)
	if err == nil {
		//check if school is existing
		school := m.School{}
		existingSchoolQuery := handler.db.Where("id = ?", newClass.SchoolId).First(&school)
		if existingSchoolQuery.RowsAffected > 0 {
			//check if class is already existing
			existingClass := m.Class{}
			existingClassQuery := handler.db.Where("section = ? and grade_level = ?", newClass.Section, newClass.GradeLevel).First(&existingClass)

			if existingClassQuery.RowsAffected == 0 {
				saveResult := handler.db.Save(&newClass)
				if (saveResult.RowsAffected > 0) {
					c.JSON(http.StatusCreated, newClass)
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