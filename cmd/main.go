package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/israelmiranda/go-cloud-run/internal/clients/viacep"
	"github.com/israelmiranda/go-cloud-run/internal/clients/weatherapi"
	"github.com/israelmiranda/go-cloud-run/internal/config"
	"github.com/israelmiranda/go-cloud-run/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("cmd/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	cfg := config.Load()

	viacepClient := viacep.NewClient(http.DefaultClient)
	weatherapiClient := weatherapi.NewClient(http.DefaultClient, cfg.WeatherApiKey)
	weatherHandler := handler.NewWeatherHandler(viacepClient, weatherapiClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", weatherHandler.HandleWeatherRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Server starting on %s...", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("\nShutdown signal received. Cleaning up...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly.")
}
