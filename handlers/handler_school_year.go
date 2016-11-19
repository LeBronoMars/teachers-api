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
	currentYear := time.Now().Year()
	sy := fmt.Sprintf("S.Y. %v - %v", currentYear, currentYear + 1)
	existingSchoolYearQuery := handler.db.Where("from = ? AND to = ?", currentYear, currentYear+1).First(&existingSchoolYear)
	if (existingSchoolYearQuery.RowsAffected > 0) {
		respond(http.StatusBadRequest, fmt.Sprintf("%v already existing", sy), c, true)
	} else {
		newSchoolYear := m.SchoolYear{}
		newSchoolYear.From = currentYear
		newSchoolYear.To = currentYear + 1
		newSchoolYear.Description = c.PostForm("description")

		saveResult := handler.db.Create(&newSchoolYear)

		if (saveResult.RowsAffected > 0) {
			respond(http.StatusCreated, fmt.Sprintf("%v successfully created", sy), c, false)
		} else {
			respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
		}
	}
	return
}