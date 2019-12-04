package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
)

// ChangePasswordParam .
type ChangePasswordParam struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required,gte=10,alphanum"`
}

// UsersResponse .
type UsersResponse struct {
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	HasNext  bool      `json:"has_next"`
	Users    []db.User `json:"users"`
}

// Users .
func Users(c *gin.Context) {
	con := DBConnection()
	defer con.Close()

	// TODO(hs.lee):
	// 검색 조건을 추가해야 한다.
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
	con.Unscoped().Order("id desc").Limit(pageSize + 1).Offset(page * pageSize).Find(&users)

	r := UsersResponse{
		Page:     page,
		PageSize: pageSize,
		HasNext:  false,
		Users:    users,
	}

	if len(users) > pageSize {
		r.HasNext = true
		r.Users = users[:len(users)-1]
	}

	c.JSON(http.StatusOK, r)
}

// User .
func User(c *gin.Context) {
	con := DBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser .
func DeleteUser(c *gin.Context) {
	con := DBConnection()
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

// ChangePassword .
func ChangePassword(c *gin.Context) {
	con := DBConnection()
	defer con.Close()

	var param ChangePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorBindJSON(err.Error()))
		return
	}

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	user.Password = param.CurrentPassword
	if !user.VerifyPassword() {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorIncorrectPassword())
		return
	}

	user.Password = param.Password
	err := user.SetPassword()
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSetPassword(err.Error()))
	}

	err = user.Save(con)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
	}

	c.Status(http.StatusOK)
}
