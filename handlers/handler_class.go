package handlers

import (
	"net/http"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
)

type ClassHandler struct {
	db *gorm.DB
}

func NewClassHandler(db *gorm.DB) *ClassHandler {
	return &ClassHandler{db}
}

func (handler ClassHandler) Index(c *gin.Context) {
	classess := []m.QryClassSchools{}		
	
	var query = handler.db

	startParam, startParamExist := c.GetQuery("start")
	limitParam, limitParamExist := c.GetQuery("limit")
	orderParam, orderParamExist := c.GetQuery("order")
	sectionParam, sectionParamExist := c.GetQuery("section")
	gradeParam, gradeParamExist := c.GetQuery("grade")
	schoolYearParam, schoolYearParamExist := c.GetQuery("school_year")

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

	//section param exist
	if sectionParamExist {
		query = query.Where("section = ?", sectionParam)
	} 

	//grade param exist
	if gradeParamExist {
		query = query.Where("grade_level = ?", gradeParam)
	} 

	//school year param exist
	if schoolYearParamExist {
		query = query.Where("school_year = ?", schoolYearParam)
	} 

	query.Where("created_by = ? AND deleted_at is NULL", GetCreator(c)).Find(&classess)
	c.JSON(http.StatusOK, classess)	
	return
}

//create new class
func (handler ClassHandler) Create(c *gin.Context) {
	var newClass m.Class
	err := c.Bind(&newClass)
	if err == nil {
		//check if school is existing
		school := m.School{}
		existingSchoolQuery := handler.db.Where("id = ?", newClass.School).First(&school)
		if existingSchoolQuery.RowsAffected > 0 {

			existingClassById := m.Class{}
			if handler.db.Where("id = ?", newClass.Id).First(&existingClassById).RowsAffected > 0 {
				existingClass := m.Class{}
				existingClassQuery := handler.db.Where("id != ? AND section = ? AND grade_level = ? AND created_by = ? AND deleted_at is NULL", newClass.Id, newClass.Section, newClass.GradeLevel, GetCreator(c)).First(&existingClass)

				if existingClassQuery.RowsAffected > 0 {
					respond(http.StatusBadRequest, fmt.Sprintf("Class with section of %s in Grade Level %v already exist.", newClass.Section, newClass.GradeLevel), c, true)
				} else {
					if (c.PostForm("for_deletion") == "false") {
						result := handler.db.Model(&existingClassById).Update(&newClass)
						if result.RowsAffected > 0 {
							updatedClass := m.Class{}
							handler.db.Where("id = ?", newClass.Id).First(&updatedClass)
							c.JSON(http.StatusOK, updatedClass)
						} else if result.Error != nil {
							respond(http.StatusBadRequest, result.Error.Error(), c, true)
						} else {
							respond(http.StatusBadRequest, "There are no changes detected.", c , true)
						}
					} else {
						if (c.PostForm("for_deletion") == "true") {
							delete := handler.db.Delete(&existingClass)
							if delete.RowsAffected > 0 {
								respond(http.StatusOK, "Record successfully deleted.", c, false)
							} else {
								respond(http.StatusBadRequest, delete.Error.Error(), c, true)
							}
						} else {
							respond(http.StatusBadRequest, "Invalid action.", c, true)
						}
					}
				}	
			} else {
				//check if class is already existing
				existingClass := m.Class{}
				existingClassQuery := handler.db.Where("section = ? AND grade_level = ? AND created_by = ? AND deleted_at is NULL", newClass.Section, newClass.GradeLevel, GetCreator(c)).First(&existingClass)

				if existingClassQuery.RowsAffected == 0 {
					newClass.CreatedBy = GetCreator(c)
					saveResult := handler.db.Create(&newClass)
					if (saveResult.RowsAffected > 0) {
						qryNewClass := m.QryClassSchools{}
						handler.db.Where("class_id = ? AND deleted_at is NULL", newClass.Id).First(&qryNewClass)
						c.JSON(http.StatusCreated, qryNewClass)
					} else {
						respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
					}
				} else {
					respond(http.StatusBadRequest, fmt.Sprintf("Class with section of %s in Grade Level %v already exist.", newClass.Section, newClass.GradeLevel), c, true)
				}
			}
		} else {
			respond(http.StatusNotFound, "School not found.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
	return
}

//show specic class
func (handler ClassHandler) Show(c *gin.Context) {
	class := m.QryClassSchools{}
	query := handler.db.Where("class_id = ? AND created_by = ? AND deleted_at is NULL", c.Param("class_id"), GetCreator(c)).First(&class)
	if query.RowsAffected > 0 {
		c.JSON(http.StatusOK, class)
	} else {
		respond(http.StatusNotFound, "Class record not found.", c, true)
	}
	return
}

//update class
func (handler ClassHandler) Update(c *gin.Context) {
	if (c.PostForm("grade_level") == "" && c.PostForm("section") == "" && c.PostForm("remarks") == "" && c.PostForm("school") == "") {
		respond(http.StatusBadRequest, "Nothing to update.", c, true)
	} else {
		class := m.Class{}
		query := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", c.Param("class_id"), GetCreator(c)).First(&class)
		if query.RowsAffected > 0 {
			//check if class is already existing
			existingClass := m.Class{}
			existingClassQuery := handler.db.Where("id != ? AND section = ? AND grade_level = ? AND created_by = ? AND deleted_at is NULL", c.Param("class_id"), c.PostForm("section"), c.PostForm("grade_level"), GetCreator(c)).First(&existingClass)
			if existingClassQuery.RowsAffected == 0 {

				if (c.PostForm("grade_level") != "") {
					class.GradeLevel = c.PostForm("grade_level")
				}

				if (c.PostForm("section") != "") {
					class.Section = c.PostForm("section")
				}

				if (c.PostForm("school_year_from") != "") {
					schoolYearFrom, _ := strconv.Atoi(c.PostForm("school_year_from"))
					if (schoolYearFrom < 2000) {
						respond(http.StatusBadRequest, "School year from must be year 2000 onwards", c, true)
						return
					} else {
						if (c.PostForm("school_year_to") != "") {
							schoolYearTo, _ := strconv.Atoi(c.PostForm("school_year_to"))
							if schoolYearFrom >= schoolYearTo {
								respond(http.StatusBadRequest, fmt.Sprintf("School year to (%v) must be greater than school year from (%v)", schoolYearTo, schoolYearFrom), c, true)
								return
							} else if schoolYearTo - schoolYearFrom > 1 {
								respond(http.StatusBadRequest, "Invalid school year duration.", c, true)
								return
							} else {
								class.SchoolYearFrom = schoolYearFrom
							}
						} else {
							if schoolYearFrom >= class.SchoolYearTo {
								respond(http.StatusBadRequest, fmt.Sprintf("School year from (%v) must be greater than school year to (%v)", schoolYearFrom, class.SchoolYearTo), c, true)
								return
							} else if class.SchoolYearTo - schoolYearFrom > 1 {
								respond(http.StatusBadRequest, "Invalid school year duration.", c, true)
								return
							} else {
								class.SchoolYearFrom = schoolYearFrom								
							}
						}
					}
				}

				if (c.PostForm("school_year_to") != "") {
					schoolYearTo, _ := strconv.Atoi(c.PostForm("school_year_to"))
					if (schoolYearTo < 2000) {
						respond(http.StatusBadRequest, "School year to must be year 2000 onwards", c, true)
						return
					} else {
						if (c.PostForm("school_year_from") != "") {
							schoolYearFrom, _ := strconv.Atoi(c.PostForm("school_year_from"))
							if schoolYearFrom >= schoolYearTo {
								respond(http.StatusBadRequest, fmt.Sprintf("School year to (%v) must be greater than school year from (%v)", schoolYearTo, schoolYearFrom), c, true)
								return
							} else if schoolYearTo - schoolYearFrom > 1 {
								respond(http.StatusBadRequest, "Invalid school year duration.", c, true)
								return
							} else {
								class.SchoolYearTo = schoolYearTo
							}
						} else {
							if class.SchoolYearFrom >= schoolYearTo {
								respond(http.StatusBadRequest, fmt.Sprintf("School year to (%v) must be greater than school year from (%v)", schoolYearTo, class.SchoolYearFrom), c, true)
								return
							} else if schoolYearTo - class.SchoolYearFrom > 1 {
								respond(http.StatusBadRequest, "Invalid school year duration.", c, true)
								return
							} else {
								class.SchoolYearTo = schoolYearTo
							} 
						}
					}
				}

				if (c.PostForm("school_id") != "") {
					school := m.School{}
					existingSchoolQuery := handler.db.Where("id = ?", c.PostForm("school_id")).First(&school)
					if existingSchoolQuery.RowsAffected == 0 {
						respond(http.StatusBadRequest, "School record not found.", c, true)
						return
					} else {
						class.School = c.PostForm("school_id")
					}
				}

				updateResult := handler.db.Save(&class)

				if updateResult.RowsAffected > 0 {
					c.JSON(http.StatusOK, class)
				} else {
					respond(http.StatusBadRequest, updateResult.Error.Error(), c, true)
				}
			} else {
				respond(http.StatusBadRequest, fmt.Sprintf("Class with section of %s in Grade Level %v already exist.", c.PostForm("grade_level"), c.PostForm("section")), c, true)
			}
		} else {
			respond(http.StatusNotFound, "Class record not found.", c, true)
		}
	}
	return
}

func (handler ClassHandler) Delete(c *gin.Context) {
	creatorId := GetCreator(c)
	class := m.Class{}
	query := handler.db.Where("id = ? AND created_by = ? AND deleted_at is NULL", c.Param("class_id"), creatorId).First(&class)
	if query.RowsAffected > 0 {
		existingClassSubject := m.ClassSubject{}
		existingClassSubjectQuery := handler.db.Where("created_by = ? AND class_id = ? AND deleted_at is NULL", creatorId, c.Param("class_id")).First(&existingClassSubject)

		if existingClassSubjectQuery.RowsAffected == 0 {
			deleteResult := handler.db.Delete(&class)
			if deleteResult.RowsAffected > 0 {
				respond(http.StatusOK, "Class record successfully deleted.", c, false)
			} else {
				respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
			}
		} else {
			respond(http.StatusBadRequest, "Unable to delete record, this class is related in class subject assignment.", c, true)
		}
	} else {
		respond(http.StatusNotFound, "Class record not found.", c, true)
	}
	return
}




