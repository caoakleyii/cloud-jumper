package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestDelayedPostPassword(t *testing.T) {
	a := New()
	a.Post("/hash", postPassword)
	a.Get("/health", GetHealth)
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	elapsed := make(chan time.Duration)
	start := time.Now()
	// make a call that should take ~5 seconds to return
	go func(t *testing.T) {

		resp, err := http.PostForm("http://localhost:9999/hash",
			url.Values{"password": {"angryMonkey"}})

		if err != nil {
			t.Errorf("TestDelayedPostPassword errored when making request to test server \n\n %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("TestDelayedPostPassword error when attempting to read body stream \n\n %v", err)
		}

		s := string(body[:])

		if s == "" {
			t.Errorf("TestDelayedPostPassword did not return a hashed password. Returned empty string")
		}
		if resp.StatusCode != 201 {
			t.Errorf("TestDelayedPostPassword did not return a 201 CREATED status. Returned Value %v", resp.StatusCode)
		}

		elapsed <- time.Since(start)
	}(t)

	// confirm that it's non-blocking
	go func(t *testing.T) {
		resp, err := http.Get("http://localhost:9999/health")

		if err != nil {
			t.Errorf("TestDelayedPostPassword error when making request to test server \n\n %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("TestDelayedPassword non-blocking health check did not return 200 OK. Returned Value %v", resp.StatusCode)
		}

		healthElapsed := time.Since(start)

		if healthElapsed > time.Duration(5*time.Second) {
			t.Errorf("TestDelayedPassword health check took longer than 5 seconds, it should be quicker and not blocked")
		}
	}(t)

	v := <-elapsed

	if v < time.Duration(5*time.Second) {
		t.Errorf("TestDelayedPassword took less than 5 seconds")
	}

	server.Shutdown(context.Background())
}

func TestDelayedSavePostPassword(t *testing.T) {
	a := New()
	a.Post("/hash", PostPassword)
	a.Get("/hash/:id", GetPassword)
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	var elapsed time.Duration
	start := time.Now()

	// make a call that should return instantly, however take ~5 seconds to save id
	resp, err := http.PostForm("http://localhost:9999/hash",
		url.Values{"password": {"angryMonkey"}})

	if err != nil {
		t.Errorf("TestDelayedSavePostPassword errored when making request to test server \n\n %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("TestDelayedSavePostPassword errored when attempting to read body stream \n\n %v", err)
	}

	id := string(body[:])

	if id == "" {
		t.Errorf("TestDelayedSavePostPassword did not return a password id. Returned empty string")
	}

	if resp.StatusCode != 201 {
		t.Errorf("TestDelayedSavePostPassword did not return a 201 CREATED status. Returned Value %v", resp.StatusCode)
	}

	// confirm that it's not available yet
	resp, err = http.Get(fmt.Sprintf("http://localhost:9999/hash/%v", id))

	if err != nil {
		t.Errorf("TestDelayedSavePostPassword errored when making request to test server \n\n %v", err)
	}

	elapsed = time.Since(start)
	if resp.StatusCode == 200 && elapsed < time.Duration(5*time.Second) {
		t.Errorf("TestDelayedSavePostPassword test was able to retrieve the password before five seconds")
	}

	// sleep 6 seconds, then confirm that it is available
	time.Sleep(6 * time.Second)
	resp, err = http.Get(fmt.Sprintf("http://localhost:9999/hash/%v", id))

	if err != nil {
		t.Errorf("TestDelayedSavePostPassword errored when making request to test server \n\n %v", err)
	}

	elapsed = time.Since(start)
	if resp.StatusCode != 200 && elapsed > time.Duration(10*time.Second) {
		t.Errorf("TestDelayedSavePostPassword test was unable to retrieve the password even after 10 seconds")
	}

	server.Shutdown(context.Background())
}

func TestGetPassword(t *testing.T) {
	expected := `ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZP ZklJz0Fd7su2A+gf7Q==`
	a := New()
	a.Get("/hash/:id", GetPassword)
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	resp, err := http.Get("http://localhost:9999/hash/abc123")

	if err != nil {
		t.Errorf("TestGetPassword errored when making request to test server \n\n %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("TestGetPassword did not return an OK status. Returned: %v", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("TestGetPassword error when attempting to read body stream \n\n %v", err)
	}

	password := string(body[:])

	if password != expected {
		t.Errorf("TestGetPassword did not return expected hashed password. Expected: %v | Returned: %v", expected, password)
	}

	server.Shutdown(context.Background())
}
