package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
	"github.com/loganstone/auth/db"
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
	Links    []Link    `json:"links"`
}

// Adjust .
func (r *UsersResponse) Adjust(pageSize int) {
	if len(r.Users) > pageSize {
		r.HasNext = true
		r.Users = r.Users[:len(r.Users)-1]
	}
}

// AttachLinks .
func (r *UsersResponse) AttachLinks(page, pageSize int, emails []string) {
	v := url.Values{}
	v.Set("page_size", strconv.Itoa(pageSize))
	for _, email := range emails {
		v.Add("email", email)
	}

	links := []Link{}
	if r.HasNext {
		v.Set("page", strconv.Itoa(page+1))
		links = append(links, Link{
			Rel:    "next",
			Method: "GET",
			Href:   fmt.Sprintf("/admin/users?%s", v.Encode()),
		})
	}

	if page > 0 {
		v.Set("page", strconv.Itoa(page-1))
		links = append(links, Link{
			Rel:    "prev",
			Method: "GET",
			Href:   fmt.Sprintf("/admin/users?%s", v.Encode()),
		})
	}
	r.Links = links
}

// Users .
func Users(c *gin.Context) {
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	page, err := Page(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrResWithErr(ErrorCodeBadPage, err))
		return
	}

	pageSize, err := PageSize(c)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrResWithErr(ErrorCodeBadPageSize, err))
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
	r.AttachLinks(page, pageSize, emails)

	c.JSON(http.StatusOK, r)
}

// User .
func User(c *gin.Context) {
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser .
func DeleteUser(c *gin.Context) {
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	user := findUserOrAbort(c, con, http.StatusNoContent)
	if user == nil {
		return
	}

	if err := user.Delete(con); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeDBTransaction, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangePassword .
func ChangePassword(c *gin.Context) {
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	var param ChangePasswordParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrResWithErr(ErrorCodeBindJSON, err))
		return
	}

	user := findUserOrAbort(c, con, http.StatusNotFound)
	if user == nil {
		return
	}

	if !user.VerifyPassword(param.CurrentPassword) {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			NewErrRes(ErrorCodeIncorrectPassword))
		return
	}

	err := user.SetPassword(param.Password)
	if err != nil {
		httpStatusCode := http.StatusInternalServerError
		errRes := NewErrResWithErr(ErrorCodeSetPassword, err)
		if errors.Is(err, db.ErrorInvalidPassword) {
			httpStatusCode = http.StatusBadRequest
			errRes = NewErrResWithErr(ErrorCodeInvalidPassword, err)
		}
		c.AbortWithStatusJSON(httpStatusCode, errRes)
		return
	}

	err = user.Save(con)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeDBTransaction, err))
		return
	}

	c.Status(http.StatusOK)
}

// RenewSession .
func RenewSession(c *gin.Context) {
	conf := configs.App()
	con := DBConnOrAbort(c)
	if con == nil {
		return
	}

	user := findUserOrAbort(c, con, http.StatusNoContent)
	if user == nil {
		return
	}

	token := utils.NewJWT(conf.SessionTokenExpire)
	sessionToken, err := token.Session(user.ID, user.Email, conf.JWTSigninKey, conf.Org)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			NewErrResWithErr(ErrorCodeSignJWT, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": sessionToken})
}
