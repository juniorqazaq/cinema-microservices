package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/email"
	"notification-service/internal/subscriber"
	"github.com/nats-io/nats.go"
)

func main() {
	// NATS Connection
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// SMTP Configuration (Use App Passwords for Gmail)
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	smtpClient := email.NewSMTPClient(smtpHost, smtpPort, smtpUser, smtpPass)

	// Sender Setup
	templateDir := "internal/templates"
	sender, err := email.NewSender(smtpClient, templateDir)
	if err != nil {
		log.Fatalf("Failed to initialize email sender: %v", err)
	}

	// Subscriber Setup
	sub := subscriber.NewNATSSubscriber(nc, sender)
	if err := sub.Start(); err != nil {
		log.Fatalf("Failed to start NATS subscriber: %v", err)
	}

	log.Println("Notification service is running and listening for events...")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down Notification service...")
}
