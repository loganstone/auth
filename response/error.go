package response

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/loganstone/auth/types"
)

func errorJSON(c echo.Context, httpStatusCode, errorCode int, message string) error {
	return c.JSON(httpStatusCode,
		types.Error{
			ErrorCode: errorCode,
			Message:   message,
		})
}

// ValidateError .
func ValidateError(c echo.Context, code int, message string) error {
	return errorJSON(c, code, types.ValidateError, message)
}

// NotFoundUser .
func NotFoundUser(c echo.Context) error {
	return c.JSON(http.StatusNotFound,
		types.Error{
			ErrorCode: types.NotFoundUser,
			Message:   "not such user",
		})
}
