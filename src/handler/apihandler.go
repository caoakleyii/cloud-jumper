/*
Package handler contains functions to handle http request.
*/
package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// APIHandler defines an http.Handler that incorporates basic routing
type APIHandler struct {
	Routes        map[RouteKey]Route
	Middleware    map[int]Middleware
	PreMiddleware map[int]Middleware
}

// Route represents a registered path, handler and named params
type Route struct {
	Path       string
	Handler    func(*Context)
	NamedParam string
}

type Middleware struct {
	Handler func(*Context)
}

// RouteKey defines a two dimensional mapping key for our route map
type RouteKey struct {
	HTTPMethod, Path string
}

// CSignal exports our channel for signaling a shutdown
// this is needed for Windows Support, you can not send a signal to an app
var CSignal = make(chan os.Signal)

// New returns a new reference to an APIHandler struct
func New() *APIHandler {
	r := make(map[RouteKey]Route)
	m := make(map[int]Middleware)
	p := make(map[int]Middleware)
	return &APIHandler{r, m, p}
}

// Our regex to check if a url matches /path/:id or /path or /path/meta
var namedParamRegex = regexp.MustCompile(`^/(\w+)/(:\w+|\w+=?)$`)

// ServerHTTP implements the http.Handler interface
// as a basic routing handler
func (a APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{w, r, make(map[string]string)}
	var handler func(*Context)

	for _, v := range a.PreMiddleware {
		v.Handler(ctx)
	}
	// check easy index first, if http method and path match a route exactly
	route, ok := a.Routes[RouteKey{r.Method, strings.ToLower(r.URL.Path)}]
	if !ok {
		for k, v := range a.Routes {
			// if the route doesn't have a named param, shortcut to the next route
			if v.NamedParam == "" {
				continue
			}

			// check the requested url with the regex path of this registered route
			m, err := regexp.MatchString(k.Path, r.URL.Path)

			if err != nil {
				handler = func(ctx *Context) { ctx.String(http.StatusInternalServerError, "Internal Server Error") }
			}

			// if match, get the param value
			if m {
				pvalue := namedParamRegex.FindStringSubmatch(r.URL.Path)[2]
				ctx.Params[v.NamedParam] = pvalue
				route = v
				ok = m
				break
			}
		}
	}

	if !ok {
		handler = func(ctx *Context) { ctx.String(http.StatusNotFound, "Not Found") }
	} else {
		handler = route.Handler
	}

	// run route handler
	handler(ctx)

	// run any middleware
	for _, v := range a.Middleware {
		v.Handler(ctx)
	}

	return
}

// Pre registers a new Pre Middleware handler
// to be called before any routes are handdled
func (a APIHandler) Pre(handler func(*Context)) {
	a.PreMiddleware[len(a.PreMiddleware)] = Middleware{handler}
}

// Use registers a new Middleware handler
// to be called after any routes are handled
func (a APIHandler) Use(handler func(*Context)) {
	a.Middleware[len(a.Middleware)] = Middleware{handler}
}

// Any registers a new GET or POST route
// for the path and handler provided
func (a APIHandler) Any(path string, handler func(*Context)) {
	a.Get(path, handler)
	a.Post(path, handler)
}

// Get registers a new GET route
// for the path and handler provided
func (a APIHandler) Get(path string, handler func(*Context)) {
	r, param := checkNamedParam(path)
	if r != "" {
		path = r
	}
	a.Routes[RouteKey{http.MethodGet, strings.ToLower(path)}] = Route{path, handler, param}
}

// Post registers a new POST route
// for the path and handler provided
func (a APIHandler) Post(path string, handler func(*Context)) {
	r, param := checkNamedParam(path)
	if r != "" {
		path = r
	}
	a.Routes[RouteKey{http.MethodPost, strings.ToLower(path)}] = Route{path, handler, param}
}

// checkNamedParam checks to see if the path provided, is using a url parameter based
// on the /path/:param syntax
func checkNamedParam(path string) (string, string) {
	match := namedParamRegex.FindStringSubmatch(path)
	if len(match) < 3 {
		return "", ""
	}
	pathprefix := match[1]
	param := match[2]
	if param != "" {
		return fmt.Sprintf("^/%v/(\\w+=?)$", pathprefix), strings.Split(param, ":")[1]
	}
	return "", ""
}
