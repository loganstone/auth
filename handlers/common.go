package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/types"
	"github.com/loganstone/auth/utils"
)

const (
	testEmailFmt = "test_%s@mail.com"
	testPassword = "password"
	letter       = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Please verify your email address.</title>
</head>

<body>
    <p>Hi. Do you want to create a new account?</p>

    <p>Help us secure your account by verifying your email address ({{ .Email }})</p>

    <p><a href="{{ .SignupURL }}">Sign Up</a></p>

    <p>If you don’t use this link within {{ .ExpireMin }} minutes, it will expire.</p>

    <p>Thanks,</p>
    <p>Your friends at {{ .Organization }}.</p>

    <p>You’re receiving this email because you recently created a new account. If this wasn’t you, please ignore this email.</p>
</body>

</html>`
)

// TestEmail .
type TestEmail struct {
	Email        string
	SignupURL    string
	ExpireMin    int
	Organization string
}

// SetUpNewTestUser .
func SetUpNewTestUser(email string, pw string) (*models.User, error) {
	con := db.Connection()
	defer con.Close()
	u := models.User{Email: email}
	u.SetPassword(pw)
	err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&u).Error
	})
	if err != nil {
		return nil, err
	}
	return &u, err
}

// SendEmailForUser 는 SendToLocalPostfix func 테스트를 위해 추가 했다.
func SendEmailForUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	user := models.User{Email: c.Param("email")}
	if con.Where(&user).First(&user).RecordNotFound() {
		return c.JSON(http.StatusNotFound,
			types.Error{
				ErrorCode: types.NotFoundUser,
				Message:   "not such user",
			})
	}

	val, err := json.Marshal(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			types.Error{
				ErrorCode: types.MarshalJSONError,
				Message:   err.Error(),
			})
	}
	signed, err := utils.Sign(val)
	testEmail := TestEmail{
		Email:        user.Email,
		SignupURL:    fmt.Sprintf("http://127.0.0.1:9900/signup/%s", signed),
		ExpireMin:    5,
		Organization: "Auth",
	}

	var tpl bytes.Buffer
	t := template.Must(template.New("letter").Parse(letter))
	err = t.Execute(&tpl, testEmail)
	if err != nil {
		log.Println("executing template:", err)
	}

	email := utils.NewEmail(
		"sys", "test@mail.com", user.Email,
		"Please verify your email address.", tpl.String())
	err = email.Send()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			types.Error{
				ErrorCode: types.SendEmailError,
				Message:   err.Error(),
			})
	}

	return c.JSON(http.StatusOK, "ok")
}
