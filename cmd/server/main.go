package main

import (
	"best-structure-example/internal/config"
	"best-structure-example/internal/handler"
	"best-structure-example/internal/kafka"
	"best-structure-example/internal/repository"
	"best-structure-example/internal/server"
	"best-structure-example/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// TODO change to logrus or zap
	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Create repository
	repo := repository.NewRepository(logger)

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logger.Fatalf("Failed to create Kafka producer: %v", err)
	}

	// Initialize Kafka consumer
	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka, logger)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Initialize service layer
	svc := service.NewService(repo, kafkaProducer, logger)

	// Initialize HTTP handlers
	handlers := handler.NewHandler(svc, logger)

	// Initialize HTTP server
	srv := server.NewServer(cfg.Server, handlers.Routes())

	// Start Kafka consumer
	go func() {
		if err := kafkaConsumer.Start(context.Background(), svc.ProcessMessage); err != nil {
			logger.Fatalf("Failed to start Kafka consumer: %v", err)
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		logger.Printf("Starting HTTP server on %s", cfg.Server.Addr)
		if err := srv.Start(); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	// Stop Kafka consumer
	if err := kafkaConsumer.Stop(); err != nil {
		logger.Fatalf("Failed to stop Kafka consumer: %v", err)
	}

	// Close Kafka producer
	if err := kafkaProducer.Close(); err != nil {
		logger.Fatalf("Failed to close Kafka producer: %v", err)
	}

	logger.Println("Server exited properly")
}
