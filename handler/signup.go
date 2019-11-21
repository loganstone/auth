package handler

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db/models"
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

    <p>Help us secure your account by verifying your email address ({{ .UserEmail }})</p>

    <p><a href="{{ .SignupURL }}">Sign Up</a></p>

    <p>If you don’t use this link within {{ .ExpireMin }} minutes, it will expire.</p>

    <p>Thanks,</p>
    <p>Your friends at {{ .Organization }}.</p>

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
	Password string `json:"password" binding:"required,gte=10,alphanum"`
}

// VerificationEmailData .
type VerificationEmailData struct {
	UserEmail    string
	SignupURL    string
	ExpireMin    int
	Organization string
}

// SendVerificationEmail .
func SendVerificationEmail(c *gin.Context) {
	conf := configs.App()
	con := GetDBConnection()
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

	token := utils.NewJWTToken(conf.SignupTokenExpire)
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

	if gin.Mode() == gin.TestMode {
		c.JSON(http.StatusOK, gin.H{"token": signupToken})
		return
	}

	// TODO(hs.lee):
	// organization, from email, signup_url  은 환경 변수로 설정하도록 수정
	var body bytes.Buffer
	data := VerificationEmailData{
		UserEmail:    param.Email,
		SignupURL:    signupToken,
		ExpireMin:    conf.SignupTokenExpire / 60,
		Organization: "Auth",
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
	conf := configs.App()
	con := GetDBConnection()
	defer con.Close()

	token := c.Param("token")
	claims, err := utils.ParseJWTSignupToken(token, conf.JWTSigninKey)
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
	con := GetDBConnection()
	defer con.Close()

	var param SignupParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	claims, err := utils.ParseJWTSignupToken(param.Token, conf.JWTSigninKey)
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

	var user models.User
	user.Email = claims.Email
	user.Password = param.Password

	errRes := createNewUser(&user)
	if errRes.ErrorCode != 0 {
		httpStatusCode := http.StatusInternalServerError
		if errRes.ErrorCode == payload.ErrorCodeUserAlreadyExists {
			httpStatusCode = http.StatusBadRequest
		}
		c.AbortWithStatusJSON(httpStatusCode, errRes)
		return
	}

	c.JSON(http.StatusCreated, user)
}
