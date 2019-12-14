package handler

import (
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

	emails := c.QueryArray("email")
	var users []db.User
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
