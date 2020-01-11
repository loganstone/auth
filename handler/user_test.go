package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loganstone/auth/db"
	"github.com/loganstone/auth/payload"
)

const (
	changedPassword = "changedPassw0rd"
)

func TestUser(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Email)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody db.JSONUser
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, user.Email, resBody.Email)
	assert.Nil(t, resBody.OTPConfirmedAt)
}

func TestUserWithNonexistentEmail(t *testing.T) {
	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	nonexistentEmail := testEmail()

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", nonexistentEmail)
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errRes payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)
	assert.Equal(t, payload.ErrorCodeNotFoundUser, errRes.ErrorCode)
}

func TestDeleteUser(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUserAsOtherUser(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	otherUser := db.User{
		Email:    testEmail(),
		Password: testPassword,
	}
	err = otherUser.Create(testDBCon)
	assert.Nil(t, err)

	setAuthJWTForTest(req, &otherUser)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteUserAsAdmin(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/admin/users/%s", user.Email)
	req, err := http.NewRequest("DELETE", uri, nil)
	assert.Nil(t, err)

	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	setAuthJWTForTest(req, admin)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestChangePassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: user.Password,
		Password:        changedPassword,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Signin
	signinReqBody := SigninParam{
		Email:    user.Email,
		Password: changedPassword,
	}
	body, err = json.Marshal(signinReqBody)
	assert.Nil(t, err)

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/signin", bytes.NewReader(body))
	defer req.Body.Close()
	assert.Nil(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resBody SiginResponse
	json.NewDecoder(w.Body).Decode(&resBody)
	assert.Equal(t, signinReqBody.Email, resBody.User.Email)
	assert.NotEqual(t, "", resBody.Token)
}

func TestChangePasswordWithIncorrectCurrentPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: "incorrectcurrentpassword",
		Password:        changedPassword,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)
	assert.Equal(t, payload.ErrorCodeIncorrectPassword, errRes.ErrorCode)
}

func TestChangePasswordWithoutPassword(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	reqBody := ChangePasswordParam{
		CurrentPassword: user.Password,
	}

	body, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/password", user.Email)
	req, err := http.NewRequest("PUT", uri, bytes.NewReader(body))
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errRes payload.ErrorCodeResponse
	json.NewDecoder(w.Body).Decode(&errRes)
	assert.Equal(t, payload.ErrorCodeBindJSON, errRes.ErrorCode)
}

func TestUsersAsAdmin(t *testing.T) {
	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	userCount := 10
	users := make([]*db.User, userCount)
	for i := 0; i < userCount; i++ {
		user, err := testUser(testDBCon)
		assert.Nil(t, err)
		users[i] = user
	}

	router := New()
	w := httptest.NewRecorder()
	uri := "/admin/users"
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, admin)
	q := req.URL.Query()

	tables := []struct {
		Page     int
		PageSize int
		UsersLen int
		HasNext  bool
	}{
		{0, 3, 3, true},
		{1, 3, 3, true},
		{2, 3, 3, true},
		{3, 3, 1, false},
	}

	for _, user := range users {
		q.Add("email", user.Email)
	}

	for _, v := range tables {
		q.Set("page", fmt.Sprintf("%d", v.Page))
		q.Set("page_size", fmt.Sprintf("%d", v.PageSize))

		req.URL.RawQuery = q.Encode()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var resBody UsersResponse
		json.NewDecoder(w.Body).Decode(&resBody)

		assert.Equal(t, v.Page, resBody.Page)
		assert.Equal(t, v.PageSize, resBody.PageSize)
		assert.Equal(t, v.HasNext, resBody.HasNext)
		assert.Equal(t, v.UsersLen, len(resBody.Users))
	}
}

func TestSearchUserAsAdmin(t *testing.T) {
	admin, err := testAdmin(testDBCon)
	assert.Nil(t, err)

	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()
	w := httptest.NewRecorder()
	uri := "/admin/users"
	req, err := http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, admin)
	q := req.URL.Query()

	expected := struct {
		Page     int
		PageSize int
		UsersLen int
		HasNext  bool
	}{
		0, 2, 1, false,
	}

	q.Add("email", user.Email)

	q.Set("page", fmt.Sprintf("%d", expected.Page))
	q.Set("page_size", fmt.Sprintf("%d", expected.PageSize))

	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resBody UsersResponse
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, expected.Page, resBody.Page)
	assert.Equal(t, expected.PageSize, resBody.PageSize)
	assert.Equal(t, expected.HasNext, resBody.HasNext)
	assert.Equal(t, expected.UsersLen, len(resBody.Users))
}

func TestRenewSession(t *testing.T) {
	user, err := testUser(testDBCon)
	assert.Nil(t, err)

	router := New()

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("/users/%s/session", user.Email)
	req, err := http.NewRequest("PUT", uri, nil)
	assert.Nil(t, err)

	setAuthJWTForTest(req, user)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resToken map[string]string
	json.NewDecoder(w.Body).Decode(&resToken)
	assert.True(t, resToken["token"] != "")

	uri = fmt.Sprintf("/users/%s", user.Email)
	req, err = http.NewRequest("GET", uri, nil)
	assert.Nil(t, err)

	sessionToken := resToken["token"]
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionToken))

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resBody db.JSONUser
	json.NewDecoder(w.Body).Decode(&resBody)

	assert.Equal(t, user.Email, resBody.Email)
	assert.Nil(t, resBody.OTPConfirmedAt)
}

func BenchmarkCreateUsersWithLoop(b *testing.B) {
	userCount := 10
	for i := 0; i < b.N; i++ {
		for i := 0; i < userCount; i++ {
			_, err := testUser(testDBCon)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func BenchmarkCreateUsersWithGoroutine(b *testing.B) {
	userCount := 10
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(userCount)
		for i := 0; i < userCount; i++ {
			go func() {
				defer wg.Done()
				_, err := testUser(testDBCon)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		wg.Wait()
	}
}

func BenchmarkCreateUsersWithGoroutineOnOneCPU(b *testing.B) {
	CPUNum := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(CPUNum)
	userCount := 10
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(userCount)
		for i := 0; i < userCount; i++ {
			go func() {
				defer wg.Done()
				_, err := testUser(testDBCon)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		wg.Wait()
	}
}
