package handlers

import (
	"net/http"
    "strings"
    "strconv"
    "net/smtp"
    "log"
    "math/rand"
    "time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "sanmateo/teachers/api/models"
	"sanmateo/teachers/api/config"
	"github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

//get all users
func (handler UserHandler) Index(c *gin.Context) {
	users := []m.User{}		
	handler.db.Find(&users)
	c.JSON(http.StatusOK,users)
	return
}

//create new user
func (handler UserHandler) Create(c *gin.Context) {
	if IsTokenValid(c) {
		var user m.User
		err := c.Bind(&user)
		if err == nil {
			existingSchool := m.School{}
			if handler.db.Where("id = ?", user.SchoolId).First(&existingSchool).RowsAffected == 0 {
				respond(http.StatusPreconditionFailed, "School record not found", c, true)
			} else {
				//check if employee no. is already in use
				existingUser := m.User{}
				if handler.db.Where("employee_no = ?", user.EmployeeNo).First(&existingUser).RowsAffected > 0 {
					respond(http.StatusPreconditionFailed, "Employee no already in use.", c, true)
				} else {
					if handler.db.Where("email = ?", user.Email).First(&existingUser).RowsAffected < 1 {
						encryptedPassword := encrypt([]byte(config.GetString("CRYPT_KEY")), user.Password)
						user.Password = encryptedPassword
						result := handler.db.Create(&user)
						if result.RowsAffected > 0 {
							token := &JWT{Token: generateJWT(user.Email)}
							c.JSON(http.StatusCreated, token)
						} else {
							respond(http.StatusBadRequest, result.Error.Error(), c , true)
						}
					} else {
						respond(http.StatusForbidden, "Email already taken", c , true)	
					}
				}
			}
		} else {
			respond(http.StatusBadRequest,err.Error(),c,true)
		}
	} else {
		respond(http.StatusForbidden, "Sorry, but your session has expired!", c, true)	
	}
	return
}

//user authentication
func (handler UserHandler) Auth(c *gin.Context) {
	if IsTokenValid(c) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		if (strings.TrimSpace(email) == "") {
			respond(http.StatusBadRequest,"Email is required",c,true)
		} else if (strings.TrimSpace(password) == "") {
			respond(http.StatusBadRequest,"Password is required",c,true)
		} else {
			//check if email already existing
			user := m.User{}	
			query := handler.db.Where("email = ?",email).Find(&user)

			if query.RowsAffected < 1 {
				respond(http.StatusBadRequest,"Account not found!",c,true)
			} else {
				decryptedPassword := decrypt([]byte(config.GetString("CRYPT_KEY")), user.Password)
				//invalid password
				if decryptedPassword != password {
					respond(http.StatusBadRequest,"Account not found!",c,true)
				} else {
					//authentication successful
					token := &JWT{Token: generateJWT(user.Email)}
					c.JSON(http.StatusCreated, token)
				}					
			}
		}
	} else {
		respond(http.StatusBadRequest,"Sorry, but your session has expired!",c,true)	
	}
}

func (handler UserHandler) ChangePassword (c *gin.Context) {
	if IsTokenValid(c) {
		email := c.PostForm("email")
		oldPassword := c.PostForm("old_password")
		newPassword := c.PostForm("new_password")
		user := m.User{}
		query := handler.db.Where("email = ?",email).Find(&user)
		
		if query.RowsAffected > 0 {
			decryptedPassword := decrypt([]byte(config.GetString("CRYPT_KEY")), user.Password)
			if decryptedPassword == oldPassword {
	 			encryptedPassword := encrypt([]byte(config.GetString("CRYPT_KEY")), newPassword)
				user.Password = encryptedPassword
				result := handler.db.Save(&user)
				if result.RowsAffected > 0 {
					respond(http.StatusOK,"Password successfully changed!",c,false)
				} else {
					respond(http.StatusBadRequest,"Unable to change password",c,true)
				}
			} else {
				respond(http.StatusBadRequest,"Invalid old password",c,true)
			}
		} else {
			respond(http.StatusBadRequest,"User not found!",c,true)
		}
	} else {
		respond(http.StatusBadRequest,"Sorry, but your session has expired!",c,true)	
	}
	return
}

func (handler UserHandler) ChangeProfilePic(c *gin.Context) {
	if IsTokenValid(c) {
		user_id, _ := strconv.Atoi(c.PostForm("user_id"))
		
		teahcher := m.Teacher{}	
		qry := handler.db.Where("id = ?", user_id).First(&teahcher)
		if qry.RowsAffected > 0 {
			teahcher.PicUrl = c.PostForm("new_pic_url")
			res := handler.db.Save(&teahcher)
			if res.RowsAffected > 0 {
				respond(http.StatusOK, teahcher.PicUrl, c, false)
			} else {
				respond(http.StatusBadRequest, res.Error.Error(), c, true)	
			}
		} else {
			respond(http.StatusBadRequest, "User not found!", c, true)
		}
	} else {
		respond(http.StatusBadRequest, "Sorry, but your session has expired!", c, true)	
	}
	return
}

func (handler UserHandler) GetUserById(c *gin.Context) {
	if IsTokenValid(c) {
		user_id := c.Param("user_id")
		user := m.User{}	
		res := handler.db.Where("id = ?",user_id).First(&user)
		if res.RowsAffected > 0 {
			c.JSON(http.StatusOK,user)
		} else {
			respond(http.StatusBadRequest,"User not found!",c,true)
		}
	} else {
		respond(http.StatusBadRequest,"Sorry, but your session has expired!",c,true)	
	}
	return
}

func (handler UserHandler) ForgotPassword(c *gin.Context) {
	email := c.PostForm("email")
	user := m.User{}
	qry := handler.db.Where("email = ?", email).First(&user)

	if qry.RowsAffected > 0 {
		from := "1sanmateo.app@gmail.com"
		pass := "sanmateo851troy"

		newPassword := RandomString(12)

  		msg := "From: " + from + "\r\n" +
           	"To: " + user.Email + "\r\n" + 
           	"MIME-Version: 1.0" +  "\r\n" +
           	"Content-type: text/html" + "\r\n" +
   			"Subject: Forgot Password Request" + "\r\n\r\n" +
			"Your new password <b>" + newPassword + "</b>. Please be sure that you'll change your password immediately." + "\r\n\r\n"

		err := smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from, []string{user.Email}, []byte(msg))

		if err != nil {
			log.Printf("smtp error: %s", err)
			return
		} else {
			encryptedPassword := encrypt([]byte(config.GetString("CRYPT_KEY")), newPassword)
			user.Password = encryptedPassword
			updateResult := handler.db.Save(&user)
			if updateResult.RowsAffected > 0 {
				respond(http.StatusOK, "Your new password was successfully sent to your email",c,false)
			} else {
				respond(http.StatusBadRequest, updateResult.Error.Error(),c,true)
			}
		}
	} else {
		respond(http.StatusBadRequest, "User record not found!",c,true)
	}
	return
}

func (handler UserHandler) GetUserInfo (c *gin.Context) {
	if c.Request.Header.Get("Authorization") != "" {
		tokenString := c.Request.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString[7 : len(tokenString)], func(token *jwt.Token) (interface{}, error) {
		    return []byte(config.GetString("TOKEN_KEY")), nil
		})
		if err != nil || !token.Valid {
			respond(http.StatusUnauthorized, err.Error(), c, true)
		} else {
			claims, _ := token.Claims.(jwt.MapClaims)
			user := m.User{}
			res := handler.db.Where("email = ?", claims["iss"]).First(&user)
			if res.RowsAffected > 0 {
				c.JSON(http.StatusOK, user)
			} else {
				respond(http.StatusUnauthorized, "User record not found", c, true)
			}
		}
	}
	return
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWZYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type JWT struct {
	Token string `json:"token"`
}
