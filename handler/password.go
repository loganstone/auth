package handler

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/utils"
)

// ChangePasswordParam .
type ChangePasswordParam struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// ResetPasswrodEmailData .
type ResetPasswrodEmailData struct {
	UserEmail    string `json:"user_email"`
	ResetURL     string `json:"reset_url"`
	ExpireMin    int    `json:"expire_min"`
	Organization string `json:"organization"`
}

// ResetPasswrodEmailResponseForTest .
type ResetPasswrodEmailResponseForTest struct {
	ResetPasswrodEmailData
	ResetPasswrodToekn string `json:"reset_password_token"`
	Subject            string `json:"subject"`
	Body               string `json:"body"`
}

// ChangePassword .
func ChangePassword(c *gin.Context) {
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	var param ChangePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrResWithErr(ErrorCodeBindJSON, err))
		return
	}

	user := findUserByEmailOrAbort(
		c.Param("email"), c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	if !user.VerifyPassword(param.CurrentPassword) {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrRes(ErrorCodeIncorrectPassword))
		return
	}

	err := user.SetPassword(param.Password)
	if err != nil {
		httpStatusCode := http.StatusInternalServerError
		errRes := NewErrResWithErr(ErrorCodeSetPassword, err)
		if errors.Is(err, db.ErrorInvalidPassword) {
			httpStatusCode = http.StatusBadRequest
			errRes = NewErrResWithErr(ErrorCodeInvalidPassword, err)
		}
		c.AbortWithStatusJSON(httpStatusCode, errRes)
		return
	}

	err = user.Save(con)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeDBTransaction, err))
		return
	}

	c.Status(http.StatusOK)
}

// SendResetPasswrodEmail .
func SendResetPasswrodEmail(c *gin.Context) {
	conf := configs.App()
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	var param SendEmailParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrResWithErr(ErrorCodeBindJSON, err))
		return
	}

	user := findUserByEmailOrAbort(
		param.Email, c, con, http.StatusBadRequest)
	if user == nil {
		return
	}
	user.PasswordResetTs = int(time.Now().Unix())

	token := utils.NewJWT(conf.ResetPasswrodTokenExpire)
	resetPasswordToken, err := token.ResetPasswrod(
		param.Email, user.PasswordResetTs, conf.JWTSigninKey, conf.Org)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeSignJWT, err))
		return
	}

	if gin.Mode() == gin.DebugMode {
		log.Println("reset password token:", resetPasswordToken)
	}

	emailTmpl, err := template.New("reset password email").Parse(param.Body)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeTmplParse, err))
		return
	}

	var body bytes.Buffer
	data := ResetPasswrodEmailData{
		UserEmail:    param.Email,
		ResetURL:     conf.ResetPasswordURL(resetPasswordToken),
		ExpireMin:    conf.ResetPasswrodTokenExpire / oneMinuteSeconds,
		Organization: conf.Org,
	}

	if err := emailTmpl.Execute(&body, data); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeTmplExecute, err))
		return
	}

	if err = utils.NewEmail(
		utils.NameFromEmail(param.Email),
		conf.SupportEmail,
		param.Email,
		param.Subject,
		body.String(),
	).Send(configs.SMTP().Addr()); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeSendEmail, err))
		return
	}

	err = user.Save(con)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeDBTransaction, err))
		return
	}

	if gin.Mode() == gin.TestMode {
		c.JSON(http.StatusOK, ResetPasswrodEmailResponseForTest{
			ResetPasswrodEmailData: data,
			ResetPasswrodToekn:     resetPasswordToken,
			Subject:                param.Subject,
			Body:                   body.String(),
		})
		return
	}

	c.Status(http.StatusOK)
}
