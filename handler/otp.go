package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// ConfirmOTPParam .
type ConfirmOTPParam struct {
	OTP string `json:"otp" binding:"required,numeric"`
}

// GenerateOTP .
func GenerateOTP(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	if user.ConfirmedOTP() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorOTPAlreadyRegistered())
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
		return tx.Save(user).Error
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

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	var param ConfirmOTPParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	if user.ConfirmedOTP() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorOTPAlreadyRegistered())
		return
	}

	// TODO(hs.lee):
	// payload 에 선행 되어야 하는 API URL 을 추가 하자
	if user.OTPSecretKey == "" {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			payload.ErrorEmptyOTPSecretKey())
		return
	}

	if !user.VerifyOTP(param.OTP) {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			payload.ErrorIncorrectOTP())
		return
	}

	user.ConfirmOTP()

	// TODO(hs.lee):
	// 백업 코드 개수와 자리를 환경 변수 처리해야 한다.
	codes := utils.DigitCodes(10, 6)
	result, err := json.Marshal(codes)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorMarshalJSON(err.Error()))
		return
	}

	err = user.SetOTPBackupCodes(result)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorResponse(
				payload.ErrorCodeSetOTPBackupCodes,
				err.Error()))
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Save(user).Error
	}); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp_backup_codes": user.OTPBackupCodes})
}

// ResetOTP .
func ResetOTP(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNoContent)
	if user == nil {
		return
	}

	user.ResetOTP()

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Save(user).Error
	}); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
