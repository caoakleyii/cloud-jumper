package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetHealth(t *testing.T) {
	expected := "Server is responding"

	a := New()
	a.Get("/health", GetHealth)
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	resp, err := http.Get("http://localhost:9999/health")

	if err != nil {
		t.Errorf("TestGetHealth errored when making request to test server \n\n %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("TestGetHealth did not return an OK status. Returned: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("TestGetHealth error when attempting to read body stream \n\n %v", err)
	}

	s := string(body[:])

	if s != expected {
		t.Errorf("TestGetHealth did not return expected response. Expected: %v | Returned: %v", expected, s)
	}

	server.Shutdown(context.Background())
}
