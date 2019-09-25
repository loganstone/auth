package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

const (
	// TODO(hs.lee): 파일로 읽도록 수정
	verificationEmailTmplText = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Please verify your email address.</title>
</head>

<body>
    <p>Hi. Do you want to create a new account?</p>

    <p>Help us secure your account by verifying your email address ({{ .user_email }})</p>

    <p><a href="{{ .signup_url }}">Sign Up</a></p>

    <p>If you don’t use this link within {{ .expire_min }} minutes, it will expire.</p>

    <p>Thanks,</p>
    <p>Your friends at {{ .organization }}.</p>

    <p>You’re receiving this email because you recently created a new account. If this wasn’t you, please ignore this email.</p>
</body>

</html>`
)

var verificationEmailTmpl = template.Must(template.New("verification email").Parse(verificationEmailTmplText))

// VerificationEmailParam .
type VerificationEmailParam struct {
	Email string `json:"email" binding:"required,email"`
}

// TokenData .
type TokenData struct {
	Email     string `json:"email"`
	ExpiredAt int64  `json:"expired_at"`
}

// SendVerificationEmail .
func SendVerificationEmail(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	var user models.User
	var param VerificationEmailParam

	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	if !con.Where("email = ?", param.Email).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.UserAlreadyExists())
		return
	}

	v, err := json.Marshal(TokenData{
		Email:     param.Email,
		ExpiredAt: time.Now().Unix() + configs.App().SignupTokenExpire,
	})
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorMarshalJSON(err.Error()))
		return
	}

	token, err := utils.Sign(v)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignToken(err.Error()))
		return
	}

	// TODO(hs.lee):
	// organization, from email, signup_url  은 환경 변수로 설정하도록 수정
	var body bytes.Buffer
	data := map[string]interface{}{
		"user_email":   param.Email,
		"signup_url":   token,
		"expire_min":   configs.App().SignupTokenExpire / 60,
		"organization": "auth",
	}

	if err := verificationEmailTmpl.Execute(&body, data); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorTmplExecute(err.Error()))
		return
	}

	if err = utils.NewEmail(
		utils.NameFromEmail(param.Email),
		"auth@email.com",
		param.Email,
		"[auth] Sign up for email address.", // TODO(hs.lee): 설정 하도록 수정
		body.String(),
	).Send(); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSendEmail(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// VerifySignupToken .
func VerifySignupToken(c *gin.Context) {
	token := c.Param("token")
	decodedToken, err := utils.Load(token)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorLoadToken(err.Error()))
		return
	}

	var tokenData TokenData
	if err := json.Unmarshal(decodedToken, &tokenData); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorUnMarshalJSON(err.Error()))
		return
	}

	if tokenData.ExpiredAt < time.Now().Unix() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorExpiredToken())
		return
	}
	c.JSON(http.StatusOK, gin.H{"email": tokenData.Email})
}
