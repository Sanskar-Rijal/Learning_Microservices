package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {
	inmemRepo := repository.NewinmemRepository()
	svc := service.NewService(inmemRepo)
	mux := http.NewServeMux()
	httpHandler := &h.HttpHandler{Service: svc}

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	// implementing gracefull shutdown here
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// creating a channel that receives os signals
	shutdown := make(chan os.Signal, 1)

	// telling go which signals to listen for
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server, %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down, due to %v signal", sig)

		// creating a timeout of 10 sec so that all request can be finished first
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// tell server to shutdown gracefully
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully  %v", err)
			server.Close()
		}
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Http server error %v", err)
	}
}
