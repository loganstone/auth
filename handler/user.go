package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/response"
)

// Users .
func Users(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	var users []models.User

	con.Find(&users)

	c.JSON(http.StatusOK, users)
}

// User .
func User(c *gin.Context) {
	con := db.Connection()
	defer con.Close()

	email := c.Param("email")
	user := models.User{Email: email}

	if con.Where(&user).First(&user).RecordNotFound() {
		c.JSON(http.StatusNotFound, response.NotFoundUser())
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
		c.JSON(http.StatusBadRequest, response.BindJSONError(err.Error()))
		return
	}

	if !con.Where(&user).First(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, response.UserAlreadyExists())
		return
	}

	if err := user.SetPassword(); err != nil {
		c.JSON(http.StatusInternalServerError,
			response.SetPasswordError(err.Error()))
		return
	}

	if err := db.DoInTransaction(con, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError,
			response.DBTransactionError(err.Error()))
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
		c.JSON(http.StatusInternalServerError,
			response.DBTransactionError(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
