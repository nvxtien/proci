package main

import (
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestCustomContextNewImpl(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://localhost:8080"))

	e.GET("/").Expect().
		Status(httptest.StatusOK).
		ContentType("application/json").
		Body().Equal("{\"message\":\"Welcome User Micro Service\"}")

	//expectedName := "iris"
	//e.POST("/set").WithFormField("name", expectedName).
	//	WithJSON("{}").Expect().
	//	Status(httptest.StatusOK).
	//	Body().Equal("set session = " + expectedName)

	//e.GET("/get").Expect().
	//	Status(httptest.StatusOK).
	//	Body().Equal(expectedName)
}
