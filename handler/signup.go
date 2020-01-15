package handler

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"text/template"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// VerificationEmailParam .
type VerificationEmailParam struct {
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// SignupParam .
type SignupParam struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// VerificationEmailData .
type VerificationEmailData struct {
	UserEmail    string `json:"user_email"`
	SignupURL    string `json:"signup_url"`
	ExpireMin    int    `json:"expire_min"`
	Organization string `json:"organization"`
}

// VerificationEmailResponseForTest .
type VerificationEmailResponseForTest struct {
	VerificationEmailData
	SignupToken string `json:"signup_token"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
}

// SendVerificationEmail .
func SendVerificationEmail(c *gin.Context) {
	conf := configs.App()
	con := DBConnection()
	defer con.Close()

	var param VerificationEmailParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	if isAbortedAsUserExist(c, con, param.Email) {
		return
	}

	token := utils.NewJWT(conf.SignupTokenExpire)
	signupToken, err := token.Signup(param.Email, conf.JWTSigninKey)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignJWTToken(err.Error()))
		return
	}

	if gin.Mode() == gin.DebugMode {
		log.Println("signup token:", signupToken)
	}

	emailTmpl, err := template.New("verification email").Parse(param.Body)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorTmplParse(err.Error()))
		return
	}

	var body bytes.Buffer
	data := VerificationEmailData{
		UserEmail:    param.Email,
		SignupURL:    conf.SignupURL(signupToken),
		ExpireMin:    conf.SignupTokenExpire / 60,
		Organization: conf.Org,
	}

	if err := emailTmpl.Execute(&body, data); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorTmplExecute(err.Error()))
		return
	}

	if err = utils.NewEmail(
		utils.NameFromEmail(param.Email),
		conf.SupportEmail,
		param.Email,
		param.Subject,
		body.String(),
	).Send(); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSendEmail(err.Error()))
		return
	}

	if gin.Mode() == gin.TestMode {
		c.JSON(http.StatusOK, VerificationEmailResponseForTest{
			VerificationEmailData: data,
			SignupToken:           signupToken,
			Subject:               param.Subject,
			Body:                  body.String(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// VerifySignupToken .
func VerifySignupToken(c *gin.Context) {
	conf := configs.App()
	con := DBConnection()
	defer con.Close()

	token := c.Param("token")
	claims, err := utils.ParseSignupJWT(token, conf.JWTSigninKey)
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok || ve.Errors != jwt.ValidationErrorExpired {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				payload.ErrorParseJWTToken(err.Error()))
			return
		}
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorExpiredToken())
		return
	}

	if isAbortedAsUserExist(c, con, claims.Email) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": claims.Email})
}

// Signup .
func Signup(c *gin.Context) {
	conf := configs.App()
	con := DBConnection()
	defer con.Close()

	var param SignupParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	claims, err := utils.ParseSignupJWT(param.Token, conf.JWTSigninKey)
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if !ok || ve.Errors != jwt.ValidationErrorExpired {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				payload.ErrorParseJWTToken(err.Error()))
			return
		}
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorExpiredToken())
		return
	}

	if isAbortedAsUserExist(c, con, claims.Email) {
		return
	}

	var user db.User
	user.Email = claims.Email
	err = user.Create(con, param.Password)
	if err != nil {
		httpStatusCode := http.StatusInternalServerError
		errRes := payload.ErrorDBTransaction(err.Error())
		if errors.Is(err, db.ErrorUserAlreadyExists) {
			httpStatusCode = http.StatusBadRequest
			errRes = payload.UserAlreadyExists()
		} else if errors.Is(err, db.ErrorInvalidPassword) {
			httpStatusCode = http.StatusBadRequest
			errRes = payload.ErrorInvalidPassword(err.Error())
		} else if errors.Is(err, db.ErrorFailSetPassword) {
			errRes = payload.ErrorSetPassword(err.Error())
		}
		c.AbortWithStatusJSON(httpStatusCode, errRes)
		return
	}

	c.JSON(http.StatusCreated, user)
}
