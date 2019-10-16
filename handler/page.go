package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/loganstone/auth/configs"
)

var (
	errPageType          = errors.New("'page' must be integer")
	errPageRange         = errors.New("'page' out of integer range")
	errPageValue         = errors.New("'page' must not be less than zero")
	errPageSizeType      = errors.New("'page_size' must be integer")
	errPageSizeRange     = errors.New("'page_size' out of integer range")
	errPageSizeValue     = errors.New("'page_size' must not be less than one")
	errOverPageSizeLimit = errors.New("'page_size' over page size limit")
)

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
	conf := configs.App()
	pageSize, err := strconv.Atoi(
		c.DefaultQuery("page_size", conf.PageSize))
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

	if conf.PageSizeLimit > 0 && conf.PageSizeLimit < pageSize {
		return 0, errOverPageSizeLimit
	}

	return pageSize, nil
}
