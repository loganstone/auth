package handler

import (
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

// ResetOTPParam .
type ResetOTPParam struct {
	BackupCode string `json:"backup_code" binding:"required,numeric"`
}

func generateOTP(con *gorm.DB, user *db.User) (string, *payload.ErrorCodeResponse) {
	user.GenerateOTPSecretKey()
	uri, err := user.OTPProvisioningURI()
	if err != nil {
		errRes := payload.ErrorResponse(
			payload.ErrorCodeOTPProvisioningURI,
			err.Error())
		return "", &errRes
	}

	err = user.Save(con)
	if err != nil {
		errRes := payload.ErrorDBTransaction(err.Error())
		return "", &errRes
	}

	return uri, nil
}

func confirmOTP(con *gorm.DB, user *db.User) *payload.ErrorCodeResponse {
	user.ConfirmOTP()

	// TODO(hs.lee):
	// 백업 코드 개수와 자리를 환경 변수 처리해야 한다.
	codes := utils.DigitCodes(10, 6)
	err := user.OTPBackupCodes.Set(codes)
	if err != nil {
		errRes := payload.ErrorResponse(
			payload.ErrorCodeSetOTPBackupCodes,
			err.Error())
		return &errRes
	}

	err = user.Save(con)
	if err != nil {
		errRes := payload.ErrorDBTransaction(err.Error())
		return &errRes
	}
	return nil
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
	uri, errRes := generateOTP(con, user)
	if errRes != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError, errRes)
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

	errRes := confirmOTP(con, user)
	if errRes != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError, errRes)
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp_backup_codes": user.OTPBackupCodes.Get()})
}

// ResetOTP .
func ResetOTP(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNoContent)
	if user == nil {
		return
	}

	if !user.ConfirmedOTP() {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	if !c.GetBool("IsAdmin") {
		var param ResetOTPParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				payload.ErrorBindJSON(err.Error()))
			return
		}

		if user.OTPBackupCodes == nil {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				payload.ErrorEmptyOTPBackupCodes(
					"empty otp backup codes. contact administrator"))
			return
		}

		if _, ok := user.OTPBackupCodes.In(param.BackupCode); !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				payload.ErrorIncorrectOTP())
			return
		}
	}

	user.ResetOTP()
	if err := user.Save(con); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
