/*
	Digivance MVC Application Framework - Unit Tests
	Base Controller Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.2.0 compatibility of controller.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in controller.go
*/

package mvcapp_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/digivance/mvcapp"
)

// TestNewApplication ensures that mvcapp.NewApplication returns the expected value
func TestNewApplication(t *testing.T) {
	app := mvcapp.NewApplication()
	if app == nil {
		t.Error("Failed to create new mvcapp Application object")
	}
}

// appTestController is used to invoke default 404 functionality and test the http server
type appTestController struct {
	*mvcapp.Controller
}

// newAppTestController is the controller creator used to return our internal testController
func newAppTestController(request *http.Request) mvcapp.IController {
	return &testController{
		Controller: mvcapp.NewBaseController(request),
	}
}

// TestApplication_Run ensures that the Application.Run method operates as expected
func TestApplication_Run(t *testing.T) {
	app := mvcapp.NewApplication()
	app.HTTPPort = 8906
	app.RouteManager.RegisterController("Home", newAppTestController)

	if app == nil || app.HTTPPort != 8906 {
		t.Fatal("Failed to create mvcapp Application object")
	}

	go func() {
		app.Run()
	}()

	res, err := http.Get("http://localhost:8906")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	defer app.Stop()

	data, err := ioutil.ReadAll(res.Body)
	actualBody := string(data)

	if len(actualBody) < 50 {
		t.Error("Did not receive expected 404 result")
		t.Log(actualBody)
	}

	if err := app.Run(); err == nil {
		t.Error("Failed to block & return error when HTTPServer is clearly in use")
	}
}
