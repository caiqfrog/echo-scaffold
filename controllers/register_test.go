package controllers

import (
	"testing"
	"github.com/labstack/echo"
	"fmt"
)

type TestController struct {
	Prefix string
}

func (self TestController) HandleTest1() (method, path string, handle echo.HandlerFunc) {
	method = "POST"
	path = "/\\/:abcd"
	handle = func(c echo.Context) error {
		return nil
	}
	return
}

func TestRegister(t *testing.T) {
	e := echo.New()

	Register(e, TestController{
		Prefix: "hhh",
	})

	fmt.Println(e.Routes())
}
