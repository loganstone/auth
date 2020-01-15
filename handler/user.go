package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
	"github.com/loganstone/auth/utils"
)

// ChangePasswordParam .
type ChangePasswordParam struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// UsersResponse .
type UsersResponse struct {
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	HasNext  bool      `json:"has_next"`
	Users    []db.User `json:"users"`
}

// Adjust .
func (r *UsersResponse) Adjust(pageSize int) {
	if len(r.Users) > pageSize {
		r.HasNext = true
		r.Users = r.Users[:len(r.Users)-1]
	}
}

// Users .
func Users(c *gin.Context) {
	con := DBConnection()
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
	emails := c.QueryArray("email")
	baseQuery := con.Unscoped()
	if len(emails) > 1 {
		baseQuery = baseQuery.Where("email IN (?)", emails)
	} else if len(emails) == 1 {
		baseQuery = baseQuery.Where("email = ?", emails[0])
	}
	baseQuery = baseQuery.Order("id desc")
	baseQuery = baseQuery.Limit(pageSize + 1).Offset(page * pageSize)
	baseQuery.Find(&users)

	r := UsersResponse{
		Page:     page,
		PageSize: pageSize,
		HasNext:  false,
		Users:    users,
	}
	r.Adjust(pageSize)

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

	if !user.VerifyPassword(param.CurrentPassword) {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			payload.ErrorIncorrectPassword())
		return
	}

	err := user.SetPassword(param.Password)
	if err != nil {
		httpStatusCode := http.StatusInternalServerError
		errRes := payload.ErrorSetPassword(err.Error())
		if errors.Is(err, db.ErrorInvalidPassword) {
			httpStatusCode = http.StatusBadRequest
			errRes = payload.ErrorInvalidPassword(err.Error())
		}
		c.AbortWithStatusJSON(httpStatusCode, errRes)
		return
	}

	err = user.Save(con)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorDBTransaction(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// RenewSession .
func RenewSession(c *gin.Context) {
	conf := configs.App()
	con := DBConnection()
	defer con.Close()

	user := findUserOrAbort(c, con, http.StatusNoContent)
	if user == nil {
		return
	}

	token := utils.NewJWT(conf.SessionTokenExpire)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			payload.ErrorSignJWTToken(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": sessionToken})
}
