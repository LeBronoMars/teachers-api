package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	h "sanmateo/teachers/api/handlers"
	m "sanmateo/teachers/api/models"
	"sanmateo/teachers/api/config"
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
	"github.com/pusher/pusher-http-go"
)

func main() {
	pusher := *InitPusher()
	db := *InitDB()
	router := gin.Default()
	LoadAPIRoutes(router, &db, &pusher)
}

func LoadAPIRoutes(r *gin.Engine, db *gorm.DB, pusher *pusher.Client) {
	publicWithBasicAuth := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		"admin@avinnovz.com" : "P@ssw0rd",
		}))
	public := r.Group("/api/v1")
	private := r.Group("/api/v1")
	private.Use(Auth(config.GetString("TOKEN_KEY")))

	//manage users
	userHandler := h.NewUserHandler(db)
	public.POST("/user", userHandler.Create)
	public.POST("/login", userHandler.Auth)
	private.GET("/users", userHandler.Index)
	private.PUT("/users/:user_id", userHandler.Update)
	private.PUT("/change_password", userHandler.ChangePassword)
	private.PUT("/change_profile_pic", userHandler.ChangeProfilePic)
	private.GET("/users/:user_id", userHandler.GetUserById)
	private.GET("/me", userHandler.GetUserInfo)
	public.POST("/forgot_password", userHandler.ForgotPassword)

	//manage school years
	schoolYearHandler := h.NewSchoolYearHandler(db)
	private.GET("/school_year", schoolYearHandler.Index)
	private.POST("/school_year", schoolYearHandler.Create)

	//manage school
	schoolHandler := h.NewSchoolHandler(db)
	publicWithBasicAuth.GET("/schools", schoolHandler.Index)
	publicWithBasicAuth.GET("/schools/:school_id", schoolHandler.Show)
	publicWithBasicAuth.POST("/school", schoolHandler.Create)
	publicWithBasicAuth.PUT("/schools/:school_id", schoolHandler.Update)

	//manage class
	classHandler := h.NewClassHandler(db)
	private.GET("/class", classHandler.Index)
	private.POST("/class", classHandler.Create)
	private.GET("/class/:class_id", classHandler.Show)
	private.PUT("/class/:class_id", classHandler.Update)
	private.DELETE("/class/:class_id", classHandler.Delete)

	//manage students
	studentHandler := h.NewStudentHandler(db)
	private.GET("/students", studentHandler.Index)
	private.POST("/students", studentHandler.Create)
	private.PUT("/students/:id", studentHandler.Update)
	private.GET("/students/:id", studentHandler.Show)
	private.DELETE("/students/:id", studentHandler.Delete)

	//manage subjects
	subjectHandler := h.NewSubjectHandler(db)
	private.GET("/subjects", subjectHandler.Index)
	private.POST("/subjects", subjectHandler.Create)
	private.PUT("/subjects/:subject_id", subjectHandler.Update)
	private.GET("/subjects/:subject_id", subjectHandler.Show)
	private.DELETE("/subjects/:subject_id", subjectHandler.Delete)

	//manage class subject teacher
	classSubjectHandler := h.NewClassSubject(db)
	private.GET("/assign/class_subject", classSubjectHandler.Index)
	private.POST("/assign/class_subject", classSubjectHandler.Create)
	private.PUT("/assign/class_subject/:class_subject_id", classSubjectHandler.Update)
	private.GET("/assign/class_subject/:class_subject_id", classSubjectHandler.Show)
	private.DELETE("/assign/class_subject/:class_subject_id", classSubjectHandler.Delete)

	//manage class student
	classStudentHandler := h.NewClassStudent(db)
	private.GET("/assign/class_student", classStudentHandler.Index)
	private.POST("/assign/class_student", classStudentHandler.Create)

	//manage schedule
	scheduleHandler := h.NewScheduleHandler(db)
	private.GET("/schedule", scheduleHandler.Index)
	private.POST("/schedule", scheduleHandler.Create)
	private.DELETE("/schedule/:id", scheduleHandler.Delete)
	private.GET("/schedule/:id", scheduleHandler.Show)

	//manage attendance
	attendanceHandler := h.NewAttendanceHandler(db)
	private.GET("/attendance", attendanceHandler.Index)
	private.POST("/attendance", attendanceHandler.Create)
	private.DELETE("/attendance/:id", attendanceHandler.Delete)
	private.GET("/attendance/:id", attendanceHandler.Show)

	//manage grade status
	gradeStatusHandler := h.NewGradeStatusHandler(db)
	private.GET("/grade_status", gradeStatusHandler.Index)
	private.POST("/grade_status", gradeStatusHandler.Create)
	private.DELETE("/grade_status/:id", gradeStatusHandler.Delete)
	private.GET("/grade_status/:id", gradeStatusHandler.Show)

	//manage class diary
	classDiaryHandler := h.NewClassDiaryHandler(db)
	private.GET("/class_diary", classDiaryHandler.Index)
	private.POST("/class_diary", classDiaryHandler.Create)
	private.DELETE("/class_diary/:id", classDiaryHandler.Delete)
	private.GET("/class_diary/:id", classDiaryHandler.Show)

	r.Run(fmt.Sprintf(":%s", "8080"))
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "" {
			tokenString := c.Request.Header.Get("Authorization")

			if (strings.Contains(tokenString, "Bearer")) {
				token, err := jwt.Parse(tokenString[7 : len(tokenString)], func(token *jwt.Token) (interface{}, error) {
				    return []byte(secret), nil
				})
				if err != nil || !token.Valid {
					response := &Response{Message: err.Error()}
					c.JSON(http.StatusUnauthorized, response)
					c.Abort()
				} 
			} else {
				response := &Response{Message: "Invalid token!"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
			}
		} else {
			response := &Response{Message: "Authorization is required"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
		}
	}
}

func InitDB() *gorm.DB {
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		config.GetString("DB_USER"), config.GetString("DB_PASS"),
		config.GetString("DB_HOST"), config.GetString("DB_PORT"),
		config.GetString("DB_NAME"))
	log.Printf("\nDatabase URL: %s\n", dbURL)

	_db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database:  %s", err))
	}
	_db.DB()
	_db.LogMode(true)
	_db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&m.User{},
																&m.School{},
																&m.Class{}, 
																&m.Student{},
																&m.Subject{},
																&m.ClassSubject{},
																&m.ClassStudent{},
																&m.Schedule{},
																&m.Attendance{},
																&m.GradeStatus{},
																&m.ClassDiary{})
	return _db
}

func InitPusher() *pusher.Client {
    client := pusher.Client{
      AppId: config.GetString("PUSHER_APP_ID"),
      Key: config.GetString("PUSHER_APP_KEY"),
      Secret: config.GetString("PUSHER_APP_SECRET"),
      Cluster: config.GetString("PUSHER_CLUSTER"),
    }
    return &client
}

func GetPort() string {
    var port = os.Getenv("PORT")
    // Set a default port if there is nothing in the environment
    if port == "" {
        port = "9000"
        fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
    }
    fmt.Println("port -----> ", port)
    return ":" + port
}

type Response struct {
	Message string `json:"message"`
}