package controllers

import (
	"github.com/labstack/echo"
	"reflect"
	"strings"
	"fmt"
	"regexp"
)

var replace = regexp.MustCompile(`[/\\]+`)

func Register(e *echo.Echo, controller interface{}) {
	controllerV := reflect.ValueOf(controller)
	echoV := reflect.ValueOf(e)
	// 获得前缀
	elemV := controllerV
	for elemV.Kind() == reflect.Ptr {
		elemV = elemV.Elem()
	}
	prefix := elemV.FieldByName("Prefix").String()
	// 调用所有可见方法
	methodLen := controllerV.NumMethod()
	for i := 0;i < methodLen;i++ {
		m := controllerV.Method(i)
		result := m.Call([]reflect.Value{})
		// 获取返回参数
		method := result[0].String() // HTTP方法名称
		url := prefix + "/" + result[1].String() // 路径
		handler := result[2] // 处理函数
		// 路径规范化
		url = replace.ReplaceAllString(url, `/`)
		urlLen := len(url)
		if urlLen > 1 && url[urlLen - 1:] == `/` {
			url = string(url[:urlLen - 1])
		}
		if 0 == urlLen || url[:1] != `/` {
			url = `/` + url
		}
		// 在echo中注册
		method = strings.ToUpper(method[:1]) + strings.ToLower(method[1:])
		echoMethod := echoV.MethodByName(method)
		echoMethod.Call([]reflect.Value{
			reflect.ValueOf(url),
			handler,
		})
		fmt.Println("register router: ", method, url)
	}
}