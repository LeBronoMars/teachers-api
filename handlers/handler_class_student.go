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
	qryClassStudent := []m.QryClassStudents{}	
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	subjectCodeParam, subjectCodeParamExist := c.GetQuery("subject_code")
	teacherEmpNoParam, teacherEmpNoParamExist := c.GetQuery("teacher_employee_no")
	gradeLevelParam, gradeLevelParamExist := c.GetQuery("class_grade_level")
	classSectionParam, classSectionParamExist := c.GetQuery("class_section")

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

	if subjectCodeParamExist {
		query = query.Where("subject_code = ?", subjectCodeParam)
	}

	if teacherEmpNoParamExist {
		query = query.Where("teacher_employee_no = ?", teacherEmpNoParam)
	}

	if gradeLevelParamExist {
		query = query.Where("class_grade_level = ?", gradeLevelParam)
	}

	if classSectionParamExist {
		query = query.Where("class_section = ?", classSectionParam)
	}

	query.Where("class_student_created_by = ? AND class_student_deleted_at is NULL", GetCreator(c)).Find(&qryClassStudent)
	//query.Where("class_subject_created_by = ?", GetCreator(c)).Find(&qryClassStudent)
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
		qrySubjectClass := m.QryClassSubjects{}
		qrySubjectClassQuery := handler.db.Where("class_subject_id = ?", classStudent.ClassSubjectId).First(&qrySubjectClass)
	
		if qrySubjectClassQuery.RowsAffected > 0 {
			//check if student is existing
			student := m.Student{}
			studentQuery := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", classStudent.StudentId, creatorId).First(&student)

			if studentQuery.RowsAffected > 0 {
				existingClassStudent := m.QryClassStudents{}
				existingClassStudentQuery := handler.db.Where("student_id = ? AND class_subject_id = ? AND class_student_created_by = ? AND class_student_deleted_at is NULL", classStudent.StudentId, classStudent.ClassSubjectId, creatorId).First(&existingClassStudent)

				if existingClassStudentQuery.RowsAffected == 0 {
					if (c.PostForm("id") == "") {
						classStudent.Id = GenerateID()
					}
					classStudent.CreatedBy = creatorId
					saveResult := handler.db.Create(&classStudent)
					if (saveResult.RowsAffected > 0) {
						qryClassStudent := m.QryClassStudents{}
						qryClassStudentQuery := handler.db.Where("class_student_id = ?", classStudent.Id).First(&qryClassStudent)
						if qryClassStudentQuery.RowsAffected > 0 {
							c.JSON(http.StatusCreated, qryClassStudent)
						} else {
							respond(http.StatusBadRequest, qryClassStudentQuery.Error.Error(), c, true)
						}
					} else {
						respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
					}
				} else {
					respond(http.StatusBadRequest, "Record already exist.", c, true)					
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
