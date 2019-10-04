package handler

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

const (
	// TODO(hs.lee): 파일로 읽도록 수정
	verificationEmailTitle = "[auth] Sign up for email address."
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

// SignupParam .
type SignupParam struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	token := utils.NewJWTToken(configs.App().SignupTokenExpire)
	signupToken, err := token.Signup(param.Email)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignJWTToken(err.Error()))
		return
	}

	if gin.Mode() == gin.DebugMode {
		log.Println("signup token:", signupToken)
	}

	if gin.Mode() == gin.TestMode {
		c.JSON(http.StatusOK, gin.H{"token": signupToken})
		return
	}

	// TODO(hs.lee):
	// organization, from email, signup_url  은 환경 변수로 설정하도록 수정
	var body bytes.Buffer
	data := map[string]interface{}{
		"user_email":   param.Email,
		"signup_url":   signupToken, // TODO(hs.lee): url 로 변경
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
		verificationEmailTitle,
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
	decodedToken, err := utils.ParseJWTToken(token)
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

	var user models.User
	con := db.Connection()
	defer con.Close()
	if !con.Where("email = ?", decodedToken["email"]).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.UserAlreadyExists())
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": decodedToken["email"]})
}

// Signup .
func Signup(c *gin.Context) {
	var param SignupParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	decodedToken, err := utils.ParseJWTToken(param.Token)
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

	var user models.User
	con := db.Connection()
	defer con.Close()
	if !con.Where("email = ?", decodedToken["email"]).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.UserAlreadyExists())
		return
	}

	user.Email = decodedToken["email"].(string)
	user.Password = param.Password

	errPayload := createNewUser(&user)
	if errPayload != nil {
		httpStatusCode := http.StatusInternalServerError
		if errPayload["error_code"] == payload.ErrorCodeUserAlreadyExists {
			httpStatusCode = http.StatusBadRequest
		}
		c.AbortWithStatusJSON(httpStatusCode, errPayload)
		return
	}

	c.JSON(http.StatusCreated, user)
}
