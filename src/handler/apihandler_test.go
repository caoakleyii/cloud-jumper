package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServeHTTP(t *testing.T) {
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
		t.Errorf("TestServerHTTP errored occured when making request to test server: \n\n %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("TestServerHTTP Response code was not OK. Test server responded with %v", resp.StatusCode)
	}

	server.Shutdown(context.Background())
}

func TestPre(t *testing.T) {
	a := New()
	a.Pre(func(ctx *Context) {
		ctx.String(200, "OK")
	})
	a.Get("/health", func(ctx *Context) {})
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	resp, err := http.Get("http://localhost:9999/health")

	if err != nil {
		t.Errorf("TestPre errored when making a request to test server: \n\n %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestPre errored occured when reading body: \n\n %v", err)
	}

	s := string(body[:])
	if s != "OK" {
		t.Errorf("TestPre failed to set body in pre-middleware call.")
	}

	server.Shutdown(context.Background())
}

func TestUse(t *testing.T) {
	a := New()
	a.Use(func(ctx *Context) {
		ctx.String(200, "OK")
	})
	a.Get("/health", func(ctx *Context) {})
	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}
	go func() {
		server.ListenAndServe()
	}()

	resp, err := http.Get("http://localhost:9999/health")

	if err != nil {
		t.Errorf("TestUse errored when making a request to test server: \n\n %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestUse errored occured when reading body: \n\n %v", err)
	}

	s := string(body[:])
	if s != "OK" {
		t.Errorf("TestUse failed to set body in middleware call.")
	}

	server.Shutdown(context.Background())
}

func TestAny(t *testing.T) {
	a := New()
	p := "/health"
	a.Any(p, func(*Context) {})

	r, ok := a.Routes[RouteKey{http.MethodGet, p}]

	if !ok {
		t.Errorf("APIHandler.Any is not properly registering %v methods", http.MethodGet)
	}

	if r.Handler == nil {
		t.Errorf("APIHandler.Any is not properly passing the handlerfunc for %v", http.MethodGet)
	}

	if r.NamedParam != "" {
		t.Errorf("APIHandler.Any is incorrectly creating a name pramater for %v, on method %v", p, http.MethodGet)
	}

	if r.Path != p {
		t.Errorf("APIHandler.Any is incorrectly setting the path for method %v, path provided: %v path set: %v", http.MethodGet, p, r.Path)
	}

	r, ok = a.Routes[RouteKey{http.MethodPost, p}]

	if !ok {
		t.Errorf("APIHandler.Any is not properly registering %v methods", http.MethodPost)
	}

	if !ok {
		t.Errorf("APIHandler.Any is not properly registering %v methods", http.MethodPost)
	}

	if r.Handler == nil {
		t.Errorf("APIHandler.Any is not properly passing the handlerfunc for %v", http.MethodPost)
	}

	if r.NamedParam != "" {
		t.Errorf("APIHandler.Any is incorrectly creating a name pramater for %v, on method %v", p, http.MethodPost)
	}

	if r.Path != p {
		t.Errorf("APIHandler.Any is incorrectly setting the path for method %v, path provided: %v path set: %v", http.MethodPost, p, r.Path)
	}
}

func TestGet(t *testing.T) {
	a := New()
	p := "/health"
	a.Get(p, func(*Context) {})

	r, ok := a.Routes[RouteKey{http.MethodGet, p}]

	if !ok {
		t.Errorf("APIHandler.Get is not properly registering methods")
	}

	if r.Handler == nil {
		t.Errorf("APIHandler.Get is not properly passing the handlerfunc")
	}

	if r.NamedParam != "" {
		t.Errorf("APIHandler.Get is incorrectly creating a name pramater for %v", p)
	}

	if r.Path != p {
		t.Errorf("APIHandler.Get is incorrectly setting the path. path provided: %v path set: %v", p, r.Path)
	}

	p = "/health/:id"
	a.Get(p, func(*Context) {})

	r, ok = a.Routes[RouteKey{http.MethodGet, `^/health/(\w+=?)$`}]

	if r.NamedParam != "id" {
		t.Errorf("APIHandler.Get is incorrectly created a named parameter for %v. NamedParam: %v, Expected: id", p, r.NamedParam)
	}

}

func TestPost(t *testing.T) {
	a := New()
	p := "/health"
	a.Post(p, func(*Context) {})

	r, ok := a.Routes[RouteKey{http.MethodPost, p}]

	if !ok {
		t.Errorf("APIHandler.Post is not properly registering methods")
	}

	if r.Handler == nil {
		t.Errorf("APIHandler.Post is not properly passing the handlerfunc")
	}

	if r.NamedParam != "" {
		t.Errorf("APIHandler.Post is incorrectly creating a name pramater for %v", p)
	}

	if r.Path != p {
		t.Errorf("APIHandler.Post is incorrectly setting the path. path provided: %v path set: %v", p, r.Path)
	}

	p = "/health/:id"
	a.Post(p, func(*Context) {})

	r, ok = a.Routes[RouteKey{http.MethodPost, `^/health/(\w+=?)$`}]

	if r.NamedParam != "id" {
		t.Errorf("APIHandler.Post is incorrectly created a named parameter for %v. NamedParam: %v, Expected: id", p, r.NamedParam)
	}
}
