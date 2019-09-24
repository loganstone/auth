package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/loganstone/auth/payload"
	"github.com/stretchr/testify/assert"
)

const (
	pageTestUri     = "/test/page"
	pageQueryFmt    = "?page=%d&page_size=%d"
	pageTestFullUrl = pageTestUri + pageQueryFmt
)

func pageTestHandler(c *gin.Context) {
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

	result := gin.H{
		"page":      page,
		"page_size": pageSize,
	}

	c.JSON(http.StatusOK, result)
}

func TestPage(t *testing.T) {
	page := 1
	page_size := 10
	router := NewTest()
	router.GET(pageTestUri, pageTestHandler)
	w := httptest.NewRecorder()
	uri := fmt.Sprintf(pageTestFullUrl, page, page_size)
	req, err := http.NewRequest("GET", uri, nil)

	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]int
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, page, resBody["page"])
	assert.Equal(t, page_size, resBody["page_size"])
}

func TestPageWithBadQueryParam(t *testing.T) {
	badPages := []string{
		"page=bad",
		"page=-1",
		"page=0&page_size=bad",
		"page=0&page_size=0",
	}
	router := NewTest()
	router.GET(pageTestUri, pageTestHandler)

	for _, v := range badPages {
		w := httptest.NewRecorder()
		uri := fmt.Sprintf(pageTestUri+"?%s", v)
		req, err := http.NewRequest("GET", uri, nil)

		assert.Equal(t, err, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}