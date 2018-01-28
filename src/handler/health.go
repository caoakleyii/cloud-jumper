package handler

import "net/http"

// GetHealth used to monitor the status of the API
// will be used in unit-testing of problems
func GetHealth(ctx *Context) {
	ctx.String(http.StatusOK, "Server is responding")
	return
}
