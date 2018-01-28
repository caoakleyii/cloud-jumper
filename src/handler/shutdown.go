package handler

import (
	"net/http"
	"os"
)

/*
	3. Graceful Shutdown

	The http.Server.Shutdown() method provides a
	graceful shutdown with the release of go 1.8
	The server is being run on a seperate goroutine
	while the app is waiting on a exit signal to update
	the blocking channel. The current context has no timeout
	and will defer the cancel for any existing child context,
	after the shutdown returns.

	This handler, handles the request for /shutdown and signals
	an interrupt or fills the channel
*/

// Shutdown is a handler that when called
// signals to gracefully shutdown the server, through an interrupt signal channel
func Shutdown(ctx *Context) {
	// only works in unix, linux and osx environments
	// on windows, the p.Signal does nothing
	// will be running in a docker, to unify our environments
	pid := os.Getpid()
	p, err := os.FindProcess(pid)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = p.Signal(os.Interrupt)

	if err.Error() == "not supported by windows" {
		CSignal <- os.Interrupt
	} else if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "Shutting Down.")
	return
}
