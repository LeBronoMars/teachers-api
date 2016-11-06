package handlers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type SchoolYearHandler struct {
	db *gorm.DB
}

func NewSchoolYearHandler(db *gorm.DB) *SchoolYearHandler {
	return &SchoolYearHandler{db}
}

//get all school year
func (handler SchoolYearHandler) Index(c *gin.Context) {
	schoolYears := []m.SchoolYear{}		
	handler.db.Find(&schoolYears)
	c.JSON(http.StatusOK, schoolYears)
	return
}

func (handler SchoolYearHandler) Create(c *gin.Context) {
	existingSchoolYear := m.SchoolYear{}
	existingSchoolYearQuery := handler.db.Where("year = ?", time.Now().Year()).First(&existingSchoolYear)
	if (existingSchoolYearQuery.RowsAffected > 0) {
		respond(http.StatusBadRequest, fmt.Sprintf("School Year (%v) already existing", time.Now().Year()), c, true)
	} else {
		newSchoolYear := m.SchoolYear{}
		newSchoolYear.Year = time.Now().Year()
		newSchoolYear.Description = c.PostForm("description")

		saveResult := handler.db.Create(&newSchoolYear)

		if (saveResult.RowsAffected > 0) {
			respond(http.StatusCreated, fmt.Sprintf("School year %v was successfully created", time.Now().Year()), c, false)
		} else {
			respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
		}
	}
	return
}