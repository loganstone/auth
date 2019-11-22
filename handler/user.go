package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
)

// Users .
func Users(c *gin.Context) {
	con := GetDBConnection()
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

	var users []db.User

	con.Limit(pageSize).Offset(page * pageSize).Find(&users)

	c.JSON(http.StatusOK, users)
}

// User .
func User(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser .
func DeleteUser(c *gin.Context) {
	con := GetDBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNoContent)
	if user == nil {
		return
	}

	if err := user.Delete(con); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
