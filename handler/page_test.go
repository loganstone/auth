package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	pageTestURI     = "/test/page"
	pageQueryFmt    = "?page=%d&page_size=%d"
	pageTestFullURL = pageTestURI + pageQueryFmt
)

func pageTestHandler(c *gin.Context) {
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

	result := gin.H{
		"page":      page,
		"page_size": pageSize,
	}

	c.JSON(http.StatusOK, result)
}

func TestPage(t *testing.T) {
	page := 1
	pageSize := 10
	router := New()
	r, _ := router.(*gin.Engine)
	r.GET(pageTestURI, pageTestHandler)
	w := httptest.NewRecorder()
	uri := fmt.Sprintf(pageTestFullURL, page, pageSize)
	req, err := http.NewRequest("GET", uri, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]int
	err = json.NewDecoder(w.Body).Decode(&resBody)
	assert.NoError(t, err)

	assert.Equal(t, page, resBody["page"])
	assert.Equal(t, pageSize, resBody["page_size"])
}

func TestPageWithBadQueryParam(t *testing.T) {
	badPages := []string{
		"page=bad",
		"page=-1",
		"page=0&page_size=bad",
		"page=0&page_size=0",
	}
	router := New()
	r, _ := router.(*gin.Engine)
	r.GET(pageTestURI, pageTestHandler)

	for _, v := range badPages {
		w := httptest.NewRecorder()
		uri := fmt.Sprintf(pageTestURI+"?%s", v)
		req, err := http.NewRequest("GET", uri, nil)
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}
