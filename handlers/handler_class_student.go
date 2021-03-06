package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type ClassStudent struct {
	db *gorm.DB
}

func NewClassStudent(db *gorm.DB) *ClassStudent {
	return &ClassStudent{db}
}

//get all class student
func (handler ClassStudent) Index(c *gin.Context) {
	//qryClassStudent := []m.QryClassStudents{}	
	qryClassStudent := []m.ClassStudent{}	

	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	//subjectCodeParam, subjectCodeParamExist := c.GetQuery("subject_code")
	//teacherEmpNoParam, teacherEmpNoParamExist := c.GetQuery("teacher_employee_no")
	//gradeLevelParam, gradeLevelParamExist := c.GetQuery("class_grade_level")
	//classSectionParam, classSectionParamExist := c.GetQuery("class_section")

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

	// if subjectCodeParamExist {
	// 	query = query.Where("subject_code = ?", subjectCodeParam)
	// }

	// if teacherEmpNoParamExist {
	// 	query = query.Where("teacher_employee_no = ?", teacherEmpNoParam)
	// }

	// if gradeLevelParamExist {
	// 	query = query.Where("class_grade_level = ?", gradeLevelParam)
	// }

	// if classSectionParamExist {
	// 	query = query.Where("class_section = ?", classSectionParam)
	// }

	//query.Where("class_student_created_by = ? AND class_student_deleted_at is NULL", GetCreator(c)).Find(&qryClassStudent)
	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&qryClassStudent)
	c.JSON(http.StatusOK, qryClassStudent)
	return
}

//create new student
func (handler ClassStudent) Create(c *gin.Context) {
	var classStudent m.ClassStudent
	err := c.Bind(&classStudent)

	if err == nil {
		creatorId := GetCreator(c)

		//check if class subject is existing
		existingClass := m.Class{}
		existingClassQuery := handler.db.Where("id = ?", classStudent.ClassId).First(&existingClass)
	
		if existingClassQuery.RowsAffected > 0 {
			//check if student is existing
			student := m.Student{}
			studentQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", classStudent.StudentId, creatorId).First(&student)

			if studentQuery.RowsAffected > 0 {
				existingClassStudent := m.ClassStudent{}
				existingClassStudentQuery := handler.db.Where("id = ? AND created_by = ?", classStudent.Id, GetCreator(c)).First(&existingClassStudent)

				if existingClassStudentQuery.RowsAffected == 0 {
					if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {
						classStudent.CreatedBy = creatorId
						saveResult := handler.db.Create(&classStudent)

						if (saveResult.RowsAffected > 0) {
							c.JSON(http.StatusCreated, classStudent)
						} else {
							deletedClassStudent := m.ClassStudent{}
							if handler.db.Unscoped().Where("id = ?", classStudent.Id).First(&deletedClassStudent).RowsAffected > 0 {
								c.JSON(http.StatusOK, deletedClassStudent)
							}
						}
					} else {
						if (c.PostForm("for_deletion") == "true") {
							existingClassStudent := m.ClassStudent{}
							if handler.db.Unscoped().Where("id = ?", classStudent.Id).First(&existingClassStudent).RowsAffected > 0 {
								c.JSON(http.StatusOK, existingClassStudent)		
							} else {
								classStudent.DeletedAt = GetDeletedAt(c)
								c.JSON(http.StatusOK, classStudent)	
							}
						} else {
							respond(http.StatusBadRequest, "Invalid action.", c, true)
						}
					}
				} else {
					if (c.PostForm("for_deletion") == "" || c.PostForm("for_deletion") == "false") {	
						result := handler.db.Model(&existingClassStudent).Update(&classStudent)
						if result.RowsAffected > 0 {
							updatedClassStudent := m.ClassStudent{}
							handler.db.Where("id = ?").First(&classStudent.Id)
							c.JSON(http.StatusOK, updatedClassStudent)
						} else if result.Error != nil {
							respond(http.StatusOK, result.Error.Error(), c, true)
						} else {
							respond(http.StatusOK, "There are no changes detected.", c , true)
						}
					} else {
						if (c.PostForm("for_deletion") == "true") {
							delete := handler.db.Delete(&existingClassStudent)
							if delete.RowsAffected > 0 {
								deletedClassStudent := m.ClassStudent{}
								if handler.db.Unscoped().Where("id = ?", classStudent.Id).First(&deletedClassStudent).RowsAffected > 0 {
									c.JSON(http.StatusOK, deletedClassStudent)
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
				respond(http.StatusNotFound, "Student record not found.", c, true)			
			}
		} else {
			respond(http.StatusNotFound, "Class subject not found.", c, true)			
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}
