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

	query.Where("created_by = ?", GetCreator(c)).Find(&subjects)
	c.JSON(http.StatusOK, subjects)
	return
}

//create new Subject
func (handler SubjectHandler) Create(c *gin.Context) {
	var newSubject m.Subject
	err := c.Bind(&newSubject)

	if err == nil {
		existingSubjectById := m.Subject{}
		if handler.db.Where("id = ?", newSubject.Id).First(&existingSubjectById).RowsAffected > 0 {
			existingSubject := m.Subject{}
			existingSubjectQuery := handler.db.Where("id != ? AND subject_code = ? AND created_by = ?", newSubject.Id, c.PostForm("subject_code"), GetCreator(c)).First(&existingSubject)
			
			if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
				if existingSubjectQuery.RowsAffected > 0 {
					respond(http.StatusBadRequest, "Subject code already existing.", c, true)
				} else {
					result := handler.db.Model(&existingSubjectById).Update(&newSubject)
					if result.RowsAffected > 0 {
						updatedSubject := m.Subject{}
						handler.db.Where("id = ?", newSubject.Id).First(&updatedSubject)
						c.JSON(http.StatusOK, updatedSubject)
					} else if result.Error != nil {
						respond(http.StatusOK, result.Error.Error(), c, true)
					} else {
						respond(http.StatusOK, "There are no changes detected.", c , true)
					}
				}
			} else {
				if (c.PostForm("for_deletion") == "true") {
					delete := handler.db.Delete(&existingSubjectById)
					if delete.RowsAffected > 0 {
						if handler.db.Unscoped().Where("id = ?", newSubject.Id).First(&existingSubjectById).RowsAffected > 0 {
							c.JSON(http.StatusOK, existingSubjectById)	
						}
					} else if delete.Error != nil {
						respond(http.StatusOK, delete.Error.Error(), c, true)
					} else {
						c.JSON(http.StatusOK, existingSubjectById)
					}
				} else {
					respond(http.StatusBadRequest, "Invalid action.", c, true)
				}
			}	
		} else {
			if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
				existingSubject := m.Subject{}
				existingSubjectQuery := handler.db.Where("subject_code = ? AND created_by = ? AND deleted_at is NULL", c.PostForm("subject_code"), GetCreator(c)).First(&existingSubject)

				if existingSubjectQuery.RowsAffected == 0 {
					if (c.PostForm("description") != "") {
						newSubject.Description = c.PostForm("description")				
					}
					newSubject.CreatedBy = GetCreator(c)
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
				if (c.PostForm("for_deletion") == "true") {
					if handler.db.Unscoped().Where("id = ?", newSubject.Id).First(&existingSubjectById).RowsAffected > 0 {
						c.JSON(http.StatusOK, existingSubjectById)		
					} else {
						newSubject.DeletedAt = GetDeletedAt(c)
						c.JSON(http.StatusOK, newSubject)	
					}
				} else {
					respond(http.StatusBadRequest, "Invalid action.", c, true)
				}
			}
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

func (handler SubjectHandler) Show(c *gin.Context) {
	subjectId := c.Param("subject_id")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("id = ? and created_by = ? AND deleted_at is NULL", subjectId, GetCreator(c)).First(&subject)

	if subjectQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, subject)
	} else {
		respond(http.StatusNotFound, "Subject record not found.", c, true)
	}
	return
}

func (handler SubjectHandler) Update(c *gin.Context) {
	subjectId := c.Param("subject_id")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", subjectId, GetCreator(c)).First(&subject)

	if subjectQuery.RowsAffected > 0 {
		existingSubject := m.Subject{}
		existingSubjectQuery := handler.db.Where("subject_code = ? AND id != ? AND created_by = ? AND deleted_at is NULL", c.PostForm("subject_code"), subject.Id, GetCreator(c)).First(&existingSubject)

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
		respond(http.StatusNotFound, "Subject record not found.", c, true)
	}
}

func (handler SubjectHandler) Delete(c *gin.Context) {
	creatorId := GetCreator(c)
	subjectId := c.Param("subject_id")
	subject := m.Subject{}
	subjectQuery := handler.db.Where("id = ? and created_by = ?", subjectId, creatorId).First(&subject)

	if subjectQuery.RowsAffected > 0 {
				existingClassSubject := m.ClassSubject{}
		existingClassSubjectQuery := handler.db.Where("created_by = ? AND subject_id = ? AND deleted_at is NULL", creatorId, subjectId).First(&existingClassSubject)

		if existingClassSubjectQuery.RowsAffected == 0 {
			deleteResult := handler.db.Delete(&subject)
			if deleteResult.RowsAffected > 0 {
				respond(http.StatusOK, "Subject successfully deleted.", c, false)
			} else {
				respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Unable to delete record, this subject is related in class subject assignment.", c, true)
		}
	} else {
		respond(http.StatusNotFound, "Subject record not found.", c, true)
	}
	return
}
