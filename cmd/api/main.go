package main

import (
	"context"
	"fmt"
	"log"
	"messaging/config"
	"messaging/internal/infrastructure/database"
	"messaging/internal/infrastructure/kafka"
	emailRepository "messaging/internal/infrastructure/repository/email"
	"messaging/internal/infrastructure/repository/event"
	smsRepository "messaging/internal/infrastructure/repository/sms"
	emailSender "messaging/internal/infrastructure/sender/email"
	smsSender "messaging/internal/infrastructure/sender/sms"
	"messaging/internal/usecase"
	emailSvc "messaging/internal/usecase/email"
	smsSvc "messaging/internal/usecase/sms"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config couldn't load...")
	}

	fmt.Println(cfg.SMS)
	fmt.Println(cfg.Email)
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", err, nil)
	}
	defer db.Close()

	eventDB, err := database.NewPostgres(cfg.EventDB)
	if err != nil {
		log.Fatal("Failed to connect to event database", err, nil)
	}
	defer eventDB.Close()

	eventRepo := event.NewEventStore(eventDB)
	smsRepo := smsRepository.NewSMSRepository(db)
	emailRepo := emailRepository.NewEmailRepository(db)

	smsSender := smsSender.NewSMSSender(cfg.SMS)
	emailSender := emailSender.NewEmailSender(cfg.Email)

	smsService := smsSvc.NewSMSService(eventRepo, smsRepo, smsSender)
	emailService := emailSvc.NewEmailService(eventRepo, emailRepo, emailSender)

	router := usecase.NewMessageRouter(smsService, emailService)
	fmt.Println("kafka config------------------------", cfg.Kafka)
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, router)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if consumeErr := consumer.Start(ctx); consumeErr != nil {
			log.Fatalf("Kafka consumer error: %v", consumeErr)
		}
	}()

	log.Println("Notification service started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Notification service shutting down...")
	// Cancel the context to stop the consumer loop
	cancel()

	// Close the Kafka connection
	if err := consumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}

	log.Println("Shutdown complete.")
}
