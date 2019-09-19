package handler

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/models"
	"github.com/loganstone/auth/payload"
)

const (
	defaultPageSize = "20"
)

var (
	errPageType      = errors.New("'page' must be integer")
	errPageRange     = errors.New("'page' out of integer range")
	errPageValue     = errors.New("'page' must not be less than zero")
	errPageSizeType  = errors.New("'page_size' must be integer")
	errPageSizeRange = errors.New("'page_size' out of integer range")
	errPageSizeValue = errors.New("'page_size' must not be less than one")
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

func bind(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("", Users)
		users.GET("/:email", User)
		users.POST("", CreateUser)
		users.DELETE("/:email", DeleteUser)
	}
	r.POST("signin", Signin)
}

func newRouter() *gin.Engine {
	router := gin.Default()

	bind(router)

	if gin.Mode() == gin.DebugMode {
		// Debug uri - /debug/pprof/
		pprof.Register(router)
	}

	return router
}

// New .
func New() *gin.Engine {
	return newRouter()
}

// NewTest .
func NewTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return newRouter()
}

// Page .
func Page(c *gin.Context) (int, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		e := err.(*strconv.NumError)
		if e.Err == strconv.ErrSyntax {
			return 0, errPageType

		} else if e.Err == strconv.ErrRange {
			return 0, errPageRange

		}

		return 0, err
	}

	if page < 0 {
		return 0, errPageValue
	}

	return page, nil
}

// PageSize .
func PageSize(c *gin.Context) (int, error) {
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", defaultPageSize))
	if err != nil {
		e := err.(*strconv.NumError)
		if e.Err == strconv.ErrSyntax {
			return 0, errPageSizeType

		} else if e.Err == strconv.ErrRange {
			return 0, errPageSizeRange

		}

		return 0, err
	}

	if pageSize < 1 {
		return 0, errPageSizeValue
	}

	return pageSize, nil
}
