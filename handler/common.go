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
	"github.com/loganstone/auth/utils"
)

const (
	testEmailFmt = "test-%s@email.com"
	testPassword = "Ok1234567!"
)

const (
	oneMinuteSeconds = 60
)

var (
	errNoAuthorizedUser    = errors.New("no 'AuthorizedUser'")
	errWrongAuthorizedUser = errors.New("'AuthorizedUser' not 'db.User' type")
	errNoDBConn            = errors.New("no db connection")
	errWrongDBConn         = errors.New("wrong db connection")
)

var (
	errNotFoundUser      = errors.New("not found user")
	errUserAlreadyExists = errors.New("user already exists")
	errIncorrectPassword = errors.New("incorrect Password")
	errExpiredToken      = errors.New("expired token")

	errOTPAlreadyRegistered = errors.New("OTP has already been registered")
	errNoOTPSecretKey       = errors.New("no OTP secert key")
	errIncorrectOTP         = errors.New("OTP is Incorrect")
	errNoOTPBackupCodes     = errors.New("no otp backup codes. contact administrator")
	errRequireVerifyOTP     = errors.New("required verify OTP")
)

var errMapByCode = map[int]error{
	ErrorCodeNotFoundUser:      errNotFoundUser,
	ErrorCodeUserAlreadyExists: errUserAlreadyExists,
	ErrorCodeIncorrectPassword: errIncorrectPassword,
	ErrorCodeExpiredToken:      errExpiredToken,

	ErrorCodeOTPAlreadyRegistered: errOTPAlreadyRegistered,
	ErrorCodeNoOTPSecretKey:       errNoOTPSecretKey,
	ErrorCodeIncorrectOTP:         errIncorrectOTP,
	ErrorCodeNoOTPBackupCodes:     errNoOTPBackupCodes,
	ErrorCodeRequireVerifyOTP:     errRequireVerifyOTP,

	ErrorCodeNoDBConn:    errNoDBConn,
	ErrorCodeWrongDBConn: errWrongDBConn,
}

// SendEmailParam .
type SendEmailParam struct {
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// Link .
type Link struct {
	Rel    string `json:"rel"`
	Method string `json:"method"`
	Href   string `json:"href"`
}

// ErrorCodeResponse .
type ErrorCodeResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Links        []Link `json:"links"`
}

// NewErrRes .
func NewErrRes(code int) ErrorCodeResponse {
	err, ok := errMapByCode[code]
	if !ok {
		message := fmt.Sprintf("undefiend error code(%d)", code)
		return NewErrResWithErr(code, errors.New(message))
	}
	return NewErrResWithErr(code, err)
}

// NewErrResWithLinks .
func NewErrResWithLinks(code int, links []Link) ErrorCodeResponse {
	err, ok := errMapByCode[code]
	if !ok {
		message := fmt.Sprintf("undefiend error code(%d)", code)
		return NewErrResWithErr(code, errors.New(message))
	}
	return ErrorCodeResponse{code, err.Error(), links}
}

// NewErrResWithErr .
func NewErrResWithErr(code int, err error) ErrorCodeResponse {
	return ErrorCodeResponse{code, err.Error(), nil}
}

// AuthorizedUser .
func AuthorizedUser(c *gin.Context) (user db.User, err error) {
	authorizedUser, ok := c.Get("AuthorizedUser")
	if !ok {
		err = errNoAuthorizedUser
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
			NewErrRes(ErrorCodeNoDBConn))
		return nil
	}
	dbCon, ok := con.(*gorm.DB)
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrRes(ErrorCodeWrongDBConn))
		return nil
	}
	return dbCon
}

func findUserByEmailOrAbort(email string, c *gin.Context, con *gorm.DB, httpStatusCode int) *db.User {
	if email == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil
	}

	if c.GetBool("RequesterIsAuthorizedUser") {
		user, err := AuthorizedUser(c)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				NewErrResWithErr(ErrorCodeAuthorizedUser, err))
			return nil
		}
		return &user
	}

	// TODO(hs.lee): 테스트 케이스를 추가한다.
	// NOTE(hs.lee): 관리자인 경우 삭제된 사용자도 검색 가능
	if c.GetBool("AuthorizedUserIsAdmin") {
		con = con.Unscoped()
	}

	user := findUserByEmail(email, con)
	if user == nil {
		c.AbortWithStatusJSON(
			httpStatusCode,
			NewErrRes(ErrorCodeNotFoundUser))
		return nil
	}
	return user
}

func isAbortedAsUserExist(c *gin.Context, con *gorm.DB, email string) bool {
	user := findUserByEmail(email, con)
	if user != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrRes(ErrorCodeUserAlreadyExists))
		return true
	}
	return false
}

func findUserByEmail(email string, con *gorm.DB) *db.User {
	user := db.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		return nil
	}
	return &user
}

func testEmail() string {
	return fmt.Sprintf(testEmailFmt, uuid.New().String())
}

func setAuthJWTForTest(req *http.Request, u *db.User) {
	conf := configs.App()
	token := utils.NewJWT(10)
	sessionToken, err := token.Session(u.ID, u.Email, conf.JWTSigninKey, conf.Org)
	if err != nil {
		log.Fatalf("failed generate session token: %s\n", err.Error())
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
