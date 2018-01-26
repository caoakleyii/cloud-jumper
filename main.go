package main

import (
	"log"
	"net/http"

	"github.com/caoakleyii/cloud-jumper/src/handler"
	"github.com/caoakleyii/cloud-jumper/src/server"
)

func main() {
	h := handler.New()

	server.UseRoutes(h)
	server.UseMiddleware(h)

	s := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	// Start Server
	go func() {
		log.Printf("Sever Started on %v", s.Addr)
		s.ListenAndServe()
	}()

	server.UseGracefulShutdown(s)
}
