package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/db/models"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "ok12345678"
)

var (
	errEmptySessionUser = errors.New("'SessionUser' empty")
	errWrongSessionUser = errors.New("'SessionUser' not 'models.User' type")
)

// GetLoginUser .
func GetLoginUser(c *gin.Context) (loginUser models.User, err error) {
	sessionUser, ok := c.Get("SessionUser")
	if !ok {
		err = errEmptySessionUser
		return
	}
	loginUser, ok = sessionUser.(models.User)
	if !ok {
		err = errWrongSessionUser
		return
	}
	return
}

// GetDBConnection .
func GetDBConnection() *gorm.DB {
	dbConf := configs.DB()
	return db.Connection(dbConf.ConnectionString(), dbConf.Echo)
}

func fundUserOrAbort(c *gin.Context, con *gorm.DB) *models.User {
	email := c.Param("email")
	user := models.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusNotFound, payload.NotFoundUser())
		return nil
	}
	return &user
}

func createNewUser(user *models.User) (errRes payload.ErrorCodeResponse) {
	con := GetDBConnection()
	defer con.Close()

	if !con.Where("email = ?", user.Email).First(user).RecordNotFound() {
		errRes = payload.UserAlreadyExists()
		return
	}

	if err := user.SetPassword(); err != nil {
		errRes = payload.ErrorSetPassword(err.Error())
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(user).Error
	}); err != nil {
		errRes = payload.ErrorDBTransaction(err.Error())
		return
	}
	return
}

func reloadUser(u *models.User) bool {
	con := GetDBConnection()
	defer con.Close()

	if con.Where("email = ?", u.Email).First(u).RecordNotFound() {
		return false
	}
	return true
}

func getTestEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func setSessionTokenInReqHeaderForTest(req *http.Request, u *models.User) {
	conf := configs.App()
	token := utils.NewJWTToken(10)
	sessionToken, err := token.Session(u.ID, u.Email, conf.JWTSigninKey)
	if err != nil {
		log.Fatalf("fail generate session token: %s\n", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))
}
