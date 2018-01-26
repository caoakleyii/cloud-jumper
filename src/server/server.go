package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/caoakleyii/cloud-jumper/src/handler"
	"github.com/caoakleyii/cloud-jumper/src/middleware"
)

func UseRoutes(h *handler.APIHandler) {

	h.Any("/shutdown", handler.Shutdown)
	h.Post("/hash", handler.PostPassword)
	h.Get("/hash/:id", handler.GetPassword)
	h.Get("/health", handler.GetHealth)
	h.Get("/stats", handler.GetStastics)
}

func UseMiddleware(h *handler.APIHandler) {
	h.Pre(middleware.PreStatistics)
	h.Use(middleware.Statistics)
}

/*
	3. Graceful Shutdown

	The http.Server.Shutdown() method provides a
	graceful shutdown with the release of go 1.8
	The server is being run on a seperate goroutine
	while the app is waiting on a exit signal to update
	the blocking channel. The current context has no timeout
	and will defer the cancel for any existing child context,
	after the shutdown returns.
*/
// UseGracefulShutdown ...
func UseGracefulShutdown(s *http.Server) {
	// create a channel, for an os signal, setup a notify
	// to listen for a kill or interrupt signal
	signal.Notify(handler.CSignal, os.Kill, os.Interrupt)

	// block current thread
	<-handler.CSignal

	// close global channel
	close(handler.CSignal)

	// kill or interrrupt signal was sent, continue with graceful shut down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
