package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ride-sharing/shared/env"
	"syscall"
	"time"
)

var httpAddr = env.GetString("HTTP_ADDR", ":8081")

func main() {
	log.Println("Starting API Gateway")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /trip/preview", enableCors(handleTripPreview))
	mux.HandleFunc("/ws/drivers", enableCors(handleDriversWebSocket))
	mux.HandleFunc("/ws/riders", enableCors(handleRidersWebSocket))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1) // default buffer of 1

	// go routine to start our server
	go func() {
		log.Printf("Server listening on  %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	// channel that receives os signals
	shutdown := make(chan os.Signal, 1)

	// tell Go which signals to listen for
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// os.interrupt is the signal used by us ctrl + c to stop the project
	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server, %v", err)
	case sig := <-shutdown:
		log.Printf("Server is shutting down, due to %v signal", sig)
		// create a timeout of 10 sec so that all request can be finished
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// tell server to shutdown gracefully
		// waits for current requests to finish (max 5 seconds)
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully %v", err)
			server.Close()
		}

	}
}
