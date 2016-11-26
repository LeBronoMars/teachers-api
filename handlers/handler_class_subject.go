package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type ClassSubject struct {
	db *gorm.DB
}

func NewClassSubject(db *gorm.DB) *ClassSubject {
	return &ClassSubject{db}
}

//get all class subject
func (handler ClassSubject) Index(c *gin.Context) {
	qrySubjectClass := []m.QryClassSubjects{}		
	
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

	query.Where("class_subject_created_by = ?", GetCreator(c)).Find(&qrySubjectClass)
	c.JSON(http.StatusOK, qrySubjectClass)
	return
}

//create new Subject
func (handler ClassSubject) Create(c *gin.Context) {
	var classSubject m.ClassSubject
	err := c.Bind(&classSubject)

	if err == nil {
		creatorId := GetCreator(c)
		//check if class is existing
		class := m.QryClassSchools{}
		classQuery := handler.db.Where("class_id = ? AND created_by = ?", classSubject.ClassId, creatorId).First(&class)

		if classQuery.RowsAffected > 0 {
			//check if subject is existing
			subject := m.Subject{}
			subjectQuery := handler.db.Where("id = ? and created_by = ?", classSubject.SubjectId, creatorId).First(&subject)

			if subjectQuery.RowsAffected > 0 {
				existingClassSubject := m.ClassSubject{}
				existingClassSubjectQuery := handler.db.Where("created_by = ? AND class_id = ? AND subject_id = ?", creatorId, classSubject.ClassId, classSubject.SubjectId).First(&existingClassSubject)

				if existingClassSubjectQuery.RowsAffected == 0 {
					if (c.PostForm("id") == "") {
						classSubject.Id = GenerateID()
					}
					classSubject.CreatedBy = creatorId
					saveResult := handler.db.Create(&classSubject)
					if saveResult.RowsAffected > 0 {
						qrySubjectClass := m.QryClassSubjects{}
						handler.db.Where("class_subject_id = ?", classSubject.Id).First(&qrySubjectClass)
						c.JSON(http.StatusCreated, qrySubjectClass)
					} else {
						respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
					}
				} else {
					respond(http.StatusBadRequest, "Record already exist.", c, true)	
				}
			} else {
				respond(http.StatusBadRequest, "Subject record does not exist.", c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Class record does not exist.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler ClassSubject) Show(c *gin.Context) {
	subjectCode := c.Param("subject_code")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("subject_code = ? AND created_by = ?", subjectCode, GetCreator(c)).First(&subject)

	if subjectQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, subject)
	} else {
		respond(http.StatusNotFound, "Subject record not found", c, true)
	}
	return
}

func (handler ClassSubject) Update(c *gin.Context) {
	subjectCode := c.Param("subject_code")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("subject_code = ? AND created_by = ?", subjectCode, GetCreator(c)).First(&subject)

	if subjectQuery.RowsAffected > 0 {
		existingSubject := m.Subject{}
		existingSubjectQuery := handler.db.Where("subject_code = ? AND id != ? AND created_by", c.PostForm("subject_code"), subject.Id, GetCreator(c)).First(&existingSubject)

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
