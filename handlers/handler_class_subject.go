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
	qrySubjectClass := []m.ClassSubject{}		
	
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

	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&qrySubjectClass)
	//query.Where("class_subject_created_by = ?", GetCreator(c)).Find(&qrySubjectClass)
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
		classQuery := handler.db.Where("class_id = ? AND created_by = ? AND deleted_at is NULL", classSubject.ClassId, creatorId).First(&class)

		if classQuery.RowsAffected > 0 {
			//check if subject is existing
			subject := m.Subject{}
			subjectQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", classSubject.SubjectId, creatorId).First(&subject)

			if subjectQuery.RowsAffected > 0 {
				existingClassSubject := m.ClassSubject{}
				existingClassSubjectQuery := handler.db.Where("created_by = ? AND id = ? AND deleted_at is NULL", creatorId, classSubject.Id).First(&existingClassSubject)

				if existingClassSubjectQuery.RowsAffected == 0 {
					if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
						classSubject.CreatedBy = creatorId
						saveResult := handler.db.Create(&classSubject)
						if saveResult.RowsAffected > 0 {
							qrySubjectClass := m.ClassSubject{}
							handler.db.Where("id = ?", classSubject.Id).First(&qrySubjectClass)
							c.JSON(http.StatusCreated, qrySubjectClass)
						} else {
							deletedClassSubject := m.ClassSubject{}
							if handler.db.Unscoped().Where("id = ?", classSubject.Id).First(&deletedClassSubject).RowsAffected > 0 {
								c.JSON(http.StatusOK, deletedClassSubject)
							}
						}
					} else {
						if (c.PostForm("for_deletion") == "true") {
							existingClassSubject := m.ClassSubject{}
							if handler.db.Unscoped().Where("id = ?", classSubject.Id).First(&existingClassSubject).RowsAffected > 0 {
								c.JSON(http.StatusOK, existingClassSubject)		
							} else {
								classSubject.DeletedAt = GetDeletedAt(c)
								c.JSON(http.StatusOK, classSubject)	
							}
						} else {
							respond(http.StatusBadRequest, "Invalid action.", c, true)
						}
					}
				} else {
					if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
						result := handler.db.Model(&existingClassSubject).Update(&classSubject)
						if result.RowsAffected > 0 {
							qrySubjectClass := m.ClassSubject{}
							handler.db.Where("id = ?", classSubject.Id).First(&qrySubjectClass)
							c.JSON(http.StatusOK, qrySubjectClass)
						} else if result.Error != nil {
							respond(http.StatusOK, result.Error.Error(), c, true)
						} else {
							respond(http.StatusOK, "There are no changes detected.", c , true)
						}
					} else {
						if (c.PostForm("for_deletion") == "true") {
							delete := handler.db.Delete(&existingClassSubject)
							if delete.RowsAffected > 0 {
								deletedClassSubject := m.ClassSubject{}
								if handler.db.Unscoped().Where("id = ?", classSubject.Id).First(&deletedClassSubject).RowsAffected > 0 {
									c.JSON(http.StatusOK, deletedClassSubject)
								}
							} else {
								respond(http.StatusBadRequest, delete.Error.Error(), c, true)
							}
						} else {
							respond(http.StatusBadRequest, "Invalid action.", c, true)
						}
					}
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

//update class subject
func (handler ClassSubject) Update(c *gin.Context) {
	creatorId := GetCreator(c)
	classSubjectId := c.Param("class_subject_id")
	existingClassSubject := m.ClassSubject{}
	existingClassSubjectQuery := handler.db.Where("id = ? AND created_by = ?", classSubjectId, creatorId).First(&existingClassSubject)

	if (existingClassSubjectQuery.RowsAffected > 0) {
		//check if class is existing
		class := m.QryClassSchools{}
		classQuery := handler.db.Where("class_id = ? AND created_by = ?", c.PostForm("class_id"), creatorId).First(&class)

		if classQuery.RowsAffected > 0 {
			//check if subject is existing
			subject := m.Subject{}
			subjectQuery := handler.db.Where("id = ? and created_by = ?", c.PostForm("subject_id"), creatorId).First(&subject)

			if subjectQuery.RowsAffected > 0 {
				if (c.PostForm("class_id") != "") {
					existingClassSubject.ClassId = c.PostForm("class_id")					
				}

				if (c.PostForm("subject_id") != "") {
					existingClassSubject.SubjectId = c.PostForm("subject_id")					
				}
				
				updateResult := handler.db.Save(&existingClassSubject)
				if updateResult.RowsAffected > 0 {
					qrySubjectClass := m.ClassSubject{}
					handler.db.Where("id = ?", classSubjectId).First(&qrySubjectClass)
					c.JSON(http.StatusOK, qrySubjectClass)
				} else {
					respond(http.StatusBadRequest, updateResult.Error.Error(), c, true)
				}
			} else {
				respond(http.StatusBadRequest, "Subject record does not exist.", c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Class record does not exist.", c, true)
		}
	} else {
		respond(http.StatusNotFound, "Record not found.", c, true)
	}
	return
}

func (handler ClassSubject) Show(c *gin.Context) {
	classSubjectId := c.Param("class_subject_id")
	qrySubjectClass := m.ClassSubject{}
	qrySubjectClassQuery := handler.db.Where("id = ?", classSubjectId).First(&qrySubjectClass)

	if qrySubjectClassQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, qrySubjectClass)
	} else {
		respond(http.StatusNotFound, "Class subject record not found", c, true)
	}
	return
}

func (handler ClassSubject) Delete(c *gin.Context) {
	classSubjectId := c.Param("class_subject_id")
	existingClassSubject := m.ClassSubject{}
	existingClassSubjectQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", classSubjectId, GetCreator(c)).First(&existingClassSubject)

	if (existingClassSubjectQuery.RowsAffected > 0) {
		deleteResult := handler.db.Delete(&existingClassSubject)
		if deleteResult.RowsAffected > 0 {
			respond(http.StatusOK, "Record successfully deleted.", c, false)
		} else {
			respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
		}
	} else {
		respond(http.StatusNotFound, "Record not found.", c, true)
	}
	return
}

