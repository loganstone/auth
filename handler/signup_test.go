package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loganstone/auth/utils"
	"github.com/stretchr/testify/assert"
)

func TestSendVerificationEmail(t *testing.T) {
	reqBody := map[string]string{
		"email": getTestEmail(),
	}
	body, err := json.Marshal(reqBody)

	assert.Equal(t, err, nil)

	router := NewTest()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/signup/email/verification", bytes.NewReader(body))
	defer req.Body.Close()

	assert.Equal(t, err, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody map[string]string
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.NotEqual(t, resBody["token"], "")
	token := resBody["token"]

	decodedToken, err := utils.Load(token)
	assert.Equal(t, err, nil)

	var tokenData TokenData
	err = json.Unmarshal(decodedToken, &tokenData)
	assert.Equal(t, err, nil)

	assert.Equal(t, reqBody["email"], tokenData.Email)
}
