package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
)

func createNewUser(user *models.User) (errPayload gin.H) {
	con := db.Connection()
	defer con.Close()

	if !con.Where(&user).First(&user).RecordNotFound() {
		errPayload = payload.UserAlreadyExists()
		return
	}

	if err := user.SetPassword(); err != nil {
		errPayload = payload.ErrorSetPassword(err.Error())
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	}); err != nil {
		errPayload = payload.ErrorDBTransaction(err.Error())
		return
	}
	return
}

// Users .
func Users(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	page, err := Page(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBadPage(err.Error()))
		return
	}

	pageSize, err := PageSize(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBadPage(err.Error()))
		return
	}

	var users []models.User

	con.Limit(pageSize).Offset(page * pageSize).Find(&users)

	c.JSON(http.StatusOK, users)
}

// User .
func User(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}

	if con.Where(&user).First(&user).RecordNotFound() {
		c.AbortWithStatusJSON(
			http.StatusNotFound, payload.NotFoundUser())
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser .
func CreateUser(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

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

// DeleteUser .
func DeleteUser(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}

	if con.Where(&user).First(&user).RecordNotFound() {
		c.Status(http.StatusNoContent)
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Delete(&user).Error
	}); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}