package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type StudentHandler struct {
	db *gorm.DB
}

func NewStudentHandler(db *gorm.DB) *StudentHandler {
	return &StudentHandler{db}
}

//get all students
func (handler StudentHandler) Index(c *gin.Context) {
	students := []m.Student{}		
	handler.db.Find(&students)
	c.JSON(http.StatusOK, students)
	return
}

//create new student
func (handler StudentHandler) Create(c *gin.Context) {
	var student m.Student
	err := c.Bind(&student)

	if err == nil {
		if (c.PostForm("id") == "") {
			student.Id = GenerateID()
		}
		existingStudent := m.Student{}
		existingStudentResult := handler.db.Where("student_no = ?", student.StudentNo).First(&existingStudent)
		if existingStudentResult.RowsAffected > 0 {
			respond(http.StatusBadRequest, "Student no. already used.", c, true)
		} else {
			saveResult := handler.db.Create(&student)
			if saveResult.RowsAffected > 0 {
				c.JSON(http.StatusCreated, student)
			} else {
				respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
			}
		}
	} else {
		respond(http.StatusBadRequest,err.Error(),c,true)
	}
	return
}