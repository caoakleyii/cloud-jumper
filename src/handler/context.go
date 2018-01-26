package handler

import (
	"encoding/json"
	"net/http"
)

// Context represents the context of the current HTTP request and
// defines a basic request and response structure for
// our handlers
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
}

// String sends a string response with the provided status code
func (ctx *Context) String(status int, message string) {
	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain;charsetUTF8")
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(message))
}

// JSON sends a json response with the provided status code
func (ctx *Context) JSON(status int, body interface{}) {
	j, err := json.Marshal(body)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write(j)
}

// Param returns the registered path parameter by name
func (ctx *Context) Param(name string) string {
	return ctx.Params[name]
}
