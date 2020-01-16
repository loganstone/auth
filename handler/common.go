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
	testPassword = "Ok1234567!"
)

var (
	errEmptyAuthorizedUser = errors.New("'AuthorizedUser' empty")
	errWrongAuthorizedUser = errors.New("'AuthorizedUser' not 'db.User' type")
)

// AuthorizedUser .
func AuthorizedUser(c *gin.Context) (user db.User, err error) {
	authorizedUser, ok := c.Get("AuthorizedUser")
	if !ok {
		err = errEmptyAuthorizedUser
		return
	}
	user, ok = authorizedUser.(db.User)
	if !ok {
		err = errWrongAuthorizedUser
		return
	}
	return
}

// DBConnOrAbort .
func DBConnOrAbort(c *gin.Context) *gorm.DB {
	con, ok := c.Get("DBConnection")
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBConnection("empty db connection"))
		return nil
	}
	dbCon, ok := con.(*gorm.DB)
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBConnection("wrong db connection"))
		return nil
	}
	return dbCon
}

func findUserOrAbort(c *gin.Context, con *gorm.DB, httpStatusCode int) *db.User {
	email := c.Param("email")
	if email == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}

	if c.GetBool("RequesterIsAuthorizedUser") {
		user, err := AuthorizedUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				payload.ErrorAuthorizedUser(err))
			return nil
		}
		return &user
	}

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

func testEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func setAuthJWTForTest(req *http.Request, u *db.User) {
	conf := configs.App()
	token := utils.NewJWT(10)
	sessionToken, err := token.Session(u.ID, u.Email, conf.JWTSigninKey)
	if err != nil {
		log.Fatalf("fail generate session token: %s\n", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))
}

func newUserForTest(con *gorm.DB, isAdmin bool) (*db.User, error) {
	email := testEmail()
	user := db.User{
		Email:   email,
		IsAdmin: isAdmin,
	}
	if err := user.Create(con, testPassword); err != nil {
		return nil, err
	}
	return &user, nil
}

func testUser(con *gorm.DB) (*db.User, error) {
	return newUserForTest(con, false)
}

func testAdmin(con *gorm.DB) (*db.User, error) {
	return newUserForTest(con, true)
}
