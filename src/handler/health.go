package handler

// GetHealth used to monitor the status of the API
// will be used in unit-testing of problems
func GetHealth(ctx *Context) {
	ctx.String(200, "Server is responding")
	return
}
