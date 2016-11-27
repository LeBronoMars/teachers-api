package handlers

import (
	"net/http"

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
