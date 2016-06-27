package test

import (
	"github.com/labstack/echo"
	"net/http"
)

func (t TestController) Get() (method, path string, handler echo.HandlerFunc) {
	method = "GET"
	path = ""
	handler = func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, This is echo-scaffold")
	}
	return
}
