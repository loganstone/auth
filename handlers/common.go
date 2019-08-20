package handlers

import (
	"encoding/json"
	"net/http"

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
)

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

// SendEmailForUser 는 SendEmail func 테스트를 위해 추가 했다.
// 추후 mail_test.go 추가 후 삭제해야 한다.
func SendEmailForUser(c echo.Context) error {
	con := db.Connection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}
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

	utils.SendMail("system", "test@mail.com", user.Email, "test", signed)

	return c.JSON(http.StatusOK, "ok")
}
