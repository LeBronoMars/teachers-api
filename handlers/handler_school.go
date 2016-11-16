package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type SchoolHandler struct {
	db *gorm.DB
}

func NewSchoolHandler(db *gorm.DB) *SchoolHandler {
	return &SchoolHandler{db}
}

//get all school
func (handler SchoolHandler) Index(c *gin.Context) {
	schools := []m.School{}		
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")

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

	query.Find(&schools)
	c.JSON(http.StatusOK, schools)
	return
}

//create new school
func (handler SchoolHandler) Create(c *gin.Context) {
	existingSchool := m.School{}
	schoolName := c.PostForm("school_name")

	existingSchoolQuery := handler.db.Where("school_name = ?", schoolName).First(&existingSchool)

	if (existingSchoolQuery.RowsAffected > 0) {
		respond(http.StatusBadRequest, "school name already existing.", c, true)
	} else {
		var newSchool m.School
		err := c.Bind(&newSchool)
		if err == nil {
			saveResult := handler.db.Create(&newSchool)
			if saveResult.RowsAffected > 0 {
				c.JSON(http.StatusCreated, newSchool)
			} else {
				respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
			}
		} else {
			respond(http.StatusBadRequest, err.Error(), c, true)
		}
	}
	return
}

func (handler SchoolHandler) Update(c *gin.Context) {
	schoolId, schoolIdErr := strconv.Atoi(c.Param("school_id"))

	if schoolIdErr == nil {
		existingSchool := m.School{}
		existingSchoolQuery := handler.db.Where("id = ?", schoolId).First(&existingSchool)
		if existingSchoolQuery.RowsAffected > 0 {
			if (c.PostForm("school_name") == "" && c.PostForm("school_address") == "" && c.PostForm("contact_no") == "" && c.PostForm("latitude") == "" && c.PostForm("longitude") == "") {
				respond(http.StatusBadRequest, "Nothing to update.", c, true)
			} else {
				otherSchool := m.School{}
				schoolName := c.PostForm("school_name")
				otherSchoolNameQuery := handler.db.Where("school_name = ? AND id != ?", schoolName, schoolId).First(&otherSchool)
				if (otherSchoolNameQuery.RowsAffected > 0) {
					respond(http.StatusBadRequest, "School name already existing.", c, true)
				} else {
					existingSchool.SchoolName = schoolName

					if (c.PostForm("school_address") != "") {
						existingSchool.SchoolAddress = c.PostForm("school_address")
					}

					if (c.PostForm("contact_no") != "") {
						existingSchool.ContactNo = c.PostForm("contact_no")
					}

					if (c.PostForm("latitude") != "") {
						lat, _ := strconv.ParseFloat(c.PostForm("latitude"), 32)
						existingSchool.Latitude = lat
					}

					if (c.PostForm("longitude") != "") {
						lon, _ := strconv.ParseFloat(c.PostForm("longitude"), 32)
						existingSchool.Longitude = lon
					}

					updateResult := handler.db.Save(&existingSchool)

					if updateResult.RowsAffected > 0 {
						c.JSON(http.StatusOK, existingSchool)
					} else {
						respond(http.StatusBadRequest, updateResult.Error.Error(), c, true)
					}
				}
			}
		} else {
			respond(http.StatusBadRequest, "School record not found.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, "Invalid school id.", c, true)
	}
	return
}


func (handler SchoolHandler) Show(c *gin.Context) {
	schoolId, schoolIdErr := strconv.Atoi(c.Param("school_id"))

	if schoolIdErr == nil {
		existingSchool := m.School{}
		existingSchoolQuery := handler.db.Where("id = ?", schoolId).First(&existingSchool)
		if existingSchoolQuery.RowsAffected > 0 {
			c.JSON(http.StatusOK, existingSchool)
		} else {
			respond(http.StatusBadRequest, "School record not found.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, "Invalid school id.", c, true)
	}
	return
}

