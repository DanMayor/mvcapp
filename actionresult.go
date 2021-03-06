/*
	Digivance MVC Application Framework
	Action Result Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the base action result functionality
*/

package mvcapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// ActionResult is a base level struct that implements the Execute
// method and provides the Data []byte member
type ActionResult struct {
	// StatusCode is the HTTP status code to write with this response. default is 200 ok
	StatusCode int

	// Headers is a key value pairs map of the names and values of headers to write with this response
	Headers map[string]string

	// Cookies is a collection of http cookie values to write with this response
	Cookies []*http.Cookie

	// Data is the raw byte array representing the payload to deliver
	Data []byte
}

// NewActionResult returns a new action result populated with the provided data
func NewActionResult(data []byte) *ActionResult {
	return &ActionResult{
		StatusCode: 200,
		Headers:    map[string]string{},
		Cookies:    []*http.Cookie{},
		Data:       data,
	}
}

// RawHTML was a patch method added late in v0.1.0 to provide the ability to pass through raw html content to be
// rendered by browser. Use this with caution, any code such as javascript should be stripped from the data before
// this method is called.
func RawHTML(data string) template.HTML {
	return template.HTML(data)
}

// NewViewResult returns a new ViewResult struct with the Data
// member set to the compiled templates requested
func NewViewResult(templates []string, model interface{}) (*ActionResult, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"RawHTML": RawHTML,
	}

	page, err := template.New("ViewTemplate").Funcs(funcMap).ParseFiles(templates...)

	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	if err = page.ExecuteTemplate(buffer, "mvcapp", model); err != nil {
		return nil, err
	}

	return NewActionResult(buffer.Bytes()), nil
}

// NewJSONResult returns a new JSONResult with the payload json encoded to Data
func NewJSONResult(payload interface{}) (*ActionResult, error) {
	data, err := json.Marshal(payload)
	if len(data) <= 0 || payload == nil {
		if err != nil {
			return nil, fmt.Errorf("Failed to create json payload: %s", err)
		}

		return nil, errors.New("Failed to create json payload")
	}

	return NewActionResult(data), nil
}

// AddHeader adds an http header key value pair combination to the result
func (result *ActionResult) AddHeader(key string, val string) error {
	result.Headers[key] = val
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("Failed to set http header: %s", err)
		}

		return err
	}

	return nil
}

// AddCookie adds the provided cookie to the result
func (result *ActionResult) AddCookie(cookie *http.Cookie) error {
	result.Cookies = append(result.Cookies, cookie)

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("Failed to add cookie value: %s", r)
		}

		return err
	}

	return nil
}

// Execute writes the header, cookies and data of this action result to the client.
func (result ActionResult) Execute(response http.ResponseWriter) error {
	for k, v := range result.Headers {
		response.Header().Set(k, v)
	}

	for _, cookie := range result.Cookies {
		http.SetCookie(response, cookie)
	}

	response.WriteHeader(result.StatusCode)
	response.Write(result.Data)

	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("Failed to execute action result: %s", err)
		}

		return err
	}

	return nil
}
