package handlers

import (
	"net/http"
	"strconv"

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

	query.Find(&students)
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

func (handler StudentHandler) Show(c *gin.Context) {
	studentNo := c.Param("student_no")
	student := m.Student{}
	studentQuery := handler.db.Where("student_no = ?", studentNo).First(&student)

	if studentQuery.RowsAffected > 0 {
		c.JSON(http.StatusOK, student)
	} else {
		respond(http.StatusNotFound, "Student record not found", c, true)
	}
	return
}

func (handler StudentHandler) Update(c *gin.Context) {
	studentNo := c.Param("student_no")
	student := m.Student{}
	studentQuery := handler.db.Where("student_no = ?", studentNo).First(&student)

	if studentQuery.RowsAffected > 0 {
		if (c.PostForm("student_no") != "") {
			otherStudent := m.Student{}
			otherStudentResult := handler.db.Where("student_no = ? AND id != ?", c.PostForm("student_no"), student.Id).First(&otherStudent)

			if otherStudentResult.RowsAffected > 0 {
				respond(http.StatusBadRequest, "Student no. already assigned to other student.", c, true)
				return
			} else {
				student.StudentNo = c.PostForm("student_no")
			}
		}

		if (c.PostForm("first_name") != "") {
			student.FirstName = c.PostForm("first_name")
		}

		if (c.PostForm("middle_name") != "") {
			student.MiddleName = c.PostForm("middle_name")
		}

		if (c.PostForm("last_name") != "") {
			student.LastName = c.PostForm("last_name")
		}

		if (c.PostForm("birth_date") != "") {
			student.BirthDate = c.PostForm("birth_date")
		}

		if (c.PostForm("address") != "") {
			student.BirthDate = c.PostForm("address")
		}	

		if (c.PostForm("gender") != "") {
			student.Gender = c.PostForm("gender")
		}	

		if (c.PostForm("status") != "") {
			student.Status = c.PostForm("status")
		}

		if (c.PostForm("remarks") != "") {
			student.Remarks = c.PostForm("remarks")
		}			

		if (c.PostForm("pic_url") != "") {
			student.PicUrl = c.PostForm("pic_url")
		}

		updateResult := handler.db.Save(&student)
		if updateResult.RowsAffected > 0 {
			c.JSON(http.StatusOK, student)
		} else {
			respond(http.StatusBadRequest, updateResult.Error.Error(), c, true)
		}
	} else {
		respond(http.StatusNotFound, "Student record not found", c, true)
	}
	return
}





