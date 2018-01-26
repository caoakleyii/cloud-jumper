package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	t.Run("TestString", testString)
	t.Run("TestParam", TestParam)
}

func TestString(t *testing.T) {

	expectedBody := "Hello Test World"
	expectedStatus := 201

	a := New()
	a.Get("/health", func(ctx *Context) {
		ctx.String(expectedStatus, expectedBody)
	})
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	resp, err := http.Get("http://localhost:9999/health")

	if err != nil {
		t.Errorf("TestString errored occured when making request to test server: \n\n %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestString errored occured when reading body: \n\n %v", err)
	}

	s := string(body[:])

	if s != expectedBody {
		t.Errorf("TestString did not properly return the string written. Expected: %v | Returned: %v", expectedBody, s)
	}
	if resp.StatusCode != expectedStatus {
		t.Errorf("TestString did not properly return the status expected. Expected: %v | Returned: %v", expectedStatus, resp.StatusCode)
	}
}

func TestParam(t *testing.T) {
	expected := "abc123"

	a := New()
	a.Get("/health/:id", func(ctx *Context) {
		p := ctx.Param("id")
		ctx.String(200, p)
	})
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	resp, err := http.Get(fmt.Sprintf("http://localhost:9999/health/%v", expected))

	if err != nil {
		t.Errorf("TestParam errored occured when making request to test server: \n\n %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestParam errored occured when reading body: \n\n %v", err)
	}

	s := string(body[:])
	if s != expected {
		t.Errorf("TestParam did not properly return the string written. Expected: %v | Returned: %v", expected, s)
	}

	server.Shutdown(context.Background())
}
