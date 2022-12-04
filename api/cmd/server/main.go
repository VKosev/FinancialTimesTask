package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vkosev/ft_api/internal/config"
	"github.com/vkosev/ft_api/internal/expression"
	"github.com/vkosev/ft_api/internal/persistance"
	"github.com/vkosev/ft_api/internal/web/routes"
)

const loggerPrefix = "ft.api | "

func main() {
	// initialize logger
	logger := log.New(os.Stdout, loggerPrefix, log.LstdFlags)

	// initialize configuration
	config, err := config.New()

	if err != nil {
		fmt.Println("Failed to load configuration file")
	}

	// initialize the api routes
	resolver := expression.NewResolver(logger, persistance.NewConnection())
	router := routes.Init(logger, resolver)

	// create a new server
	s := http.Server{
		Addr:         ":" + config.Server.Port, // configure the bind address
		Handler:      router,                   // set the default handler
		ErrorLog:     logger,                   // set the logger for the server
		ReadTimeout:  5 * time.Second,          // max time to read request from the client
		WriteTimeout: 10 * time.Second,         // max time to write response to the client
		IdleTimeout:  3600 * time.Second,       // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logger.Printf("Server started on port: %s", config.Server.Port)

		err := s.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting the server: %s\n", err)
			os.Exit(1)
		}

		logger.Println("Server started.")
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
