package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type SubjectHandler struct {
	db *gorm.DB
}

func NewSubjectHandler(db *gorm.DB) *SubjectHandler {
	return &SubjectHandler{db}
}

//get all subject
func (handler SubjectHandler) Index(c *gin.Context) {
	subjects := []m.Subject{}		
	
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

	query.Find(&subjects)
	c.JSON(http.StatusOK, subjects)
	return
}

//create new Subject
func (handler SubjectHandler) Create(c *gin.Context) {
	var newSubject m.Subject
	err := c.Bind(&newSubject)

	if err == nil {
		existingSubject := m.Subject{}
		existingSubjectQuery := handler.db.Where("subject_code = ?", c.PostForm("subject_code")).First(&existingSubject)

		if existingSubjectQuery.RowsAffected == 0 {
			if (c.PostForm("id") == "") {
				newSubject.Id = GenerateID()
			}
			saveResult := handler.db.Create(&newSubject)
			if saveResult.RowsAffected > 0 {
				c.JSON(http.StatusCreated, newSubject)
			} else {
				respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Subject code already existing.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler SubjectHandler) Show(c *gin.Context) {
	subjectCode := c.Param("subject_code")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("subject_code = ?", subjectCode).First(&subject)

	if subjectQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, subject)
	} else {
		respond(http.StatusNotFound, "Subject record not found", c, true)
	}
	return
}

func (handler SubjectHandler) Update(c *gin.Context) {
	subjectCode := c.Param("subject_code")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("subject_code = ?", subjectCode).First(&subject)

	if subjectQuery.RowsAffected > 0 {
		existingSubject := m.Subject{}
		existingSubjectQuery := handler.db.Where("subject_code = ? AND id != ?", c.PostForm("subject_code"), subject.Id).First(&existingSubject)

		if existingSubjectQuery.RowsAffected == 0 {
			if (c.PostForm("subject_name") != "") {
				subject.SubjectName = c.PostForm("subject_name")
			}
			
			if (c.PostForm("subject_code") != "") {
				subject.SubjectCode = c.PostForm("subject_code")
			}

			if (c.PostForm("description") != "") {
				subject.Description = c.PostForm("description")
			}

			updateResult := handler.db.Save(&subject)
			if updateResult.RowsAffected > 0 {
				c.JSON(http.StatusOK, subject)
			} else {
				respond(http.StatusBadRequest, updateResult.Error.Error(), c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Subject code already existing.", c, true)
		}
	} else {
		respond(http.StatusNotFound, "Subject record not found", c, true)
	}
}