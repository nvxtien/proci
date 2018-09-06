package main

import (
	"testing"

	"github.com/kataras/iris/httptest"
	"github.com/proci/model"
	"github.com/proci"
	"encoding/json"
	"fmt"
)

func TestCustomContextNewImpl(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app, httptest.URL("http://localhost:8080"))

	e.GET("/").Expect().
		Status(httptest.StatusOK).
		ContentType("application/json").
		Body().Equal("{\"message\":\"Welcome User Micro Service\"}")

	//e.POST("/network/create").
	//	WithJSON("{\"numberOfOrder\":2, \"\":\"\"}")
	//expectedName := "iris"
	//e.POST("/set").WithFormField("name", expectedName).
	//	WithJSON("{}").Expect().
	//	Status(httptest.StatusOK).
	//	Body().Equal("set session = " + expectedName)

	//e.GET("/get").Expect().
	//	Status(httptest.StatusOK).
	//	Body().Equal(expectedName)
}

func TestParse(t *testing.T) {
	config := &model.Configuration{}
	config.MumberOfOrg = 2
	config.OrdererType = proci.Solo
	config.Profile = "test"

	c, _ := json.Marshal(config)
	fmt.Printf("config: %s\n", string(c))

	cfg := &model.Configuration{}
	json.Unmarshal([]byte("{\"numberOfOrg\":2," +
		"\"ordererType\":\"kafka1\"," +
		"\"numberOfKafka\":0," +
		"\"profile\":\"test\"}"), cfg)
	fmt.Printf("unmarshal config: %s\n", cfg.OrdererType)

}
