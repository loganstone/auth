package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/db/models"
	"github.com/loganstone/auth/payload"
)

// GenerateOTP .
func GenerateOTP(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusNotFound, payload.NotFoundUser())
		return
	}

	user.GenerateOTPSecretKey()
	uri, err := user.OTPProvisioningURI()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorResponse(
				payload.ErrorCodeOTPProvisioningURI,
				err.Error()))
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return con.Save(&user).Error
	}); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secert_key": user.OTPSecretKey,
		"key_uri":    uri,
	})
}

// ConfirmOTP .
func ConfirmOTP(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusNotFound, payload.NotFoundUser())
		return
	}

	user.ConfirmOTP()

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return con.Save(&user).Error
	}); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	// TODO(hs.lee):
	// OTP 백업 코드를 반환해야 한다.
	c.JSON(http.StatusOK, user)
}

// ResetOTP .
func ResetOTP(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}
	if con.Where(&user).First(&user).RecordNotFound() {
		c.Status(http.StatusNoContent)
		return
	}

	user.ResetOTP()

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return con.Save(&user).Error
	}); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
