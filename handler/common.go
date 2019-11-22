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
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "ok12345678"
)

var (
	errEmptySessionUser = errors.New("'SessionUser' empty")
	errWrongSessionUser = errors.New("'SessionUser' not 'db.User' type")
)

// GetLoginUser .
func GetLoginUser(c *gin.Context) (loginUser db.User, err error) {
	sessionUser, ok := c.Get("SessionUser")
	if !ok {
		err = errEmptySessionUser
		return
	}
	loginUser, ok = sessionUser.(db.User)
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

func findUserOrAbort(c *gin.Context, con *gorm.DB, httpStatusCode int) *db.User {
	email := c.Param("email")
	user := db.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			httpStatusCode, payload.NotFoundUser())
		return nil
	}
	return &user
}

func isAbortedAsUserExist(c *gin.Context, con *gorm.DB, email string) bool {
	var count int
	con.Where("email = ?", email).Count(&count)
	if count > 1 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.UserAlreadyExists())
		return true
	}
	return false
}

func getTestEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func setSessionTokenInReqHeaderForTest(req *http.Request, u *db.User) {
	conf := configs.App()
	token := utils.NewJWTToken(10)
	sessionToken, err := token.Session(u.ID, u.Email, conf.JWTSigninKey)
	if err != nil {
		log.Fatalf("fail generate session token: %s\n", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))
}
