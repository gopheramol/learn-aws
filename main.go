package main

import (
	"context"
	"log"
	"os"

	"github.com/gopheramol/learn-aws/client"
	"github.com/gopheramol/learn-aws/service"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Environment variables loaded successfully from .env file")

	// Load AWS configuration from environment variables
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	log.Println("AWS configuration loaded successfully")

	_, credErr := cfg.Credentials.Retrieve(context.TODO())
	if credErr != nil {
		log.Fatalf("Failed to retrieve credentials, %v", credErr)
	}

	// Initialize the AWS client
	awsClient := client.NewAWSClient(cfg)
	if awsClient == nil {
		log.Fatalf("Failed to initialize AWS client")
	}
	log.Println("AWS client initialized successfully")

	// Initialize the message service
	messageService := service.NewMessageService(awsClient)
	if messageService == nil {
		log.Fatalf("Failed to initialize message service")
	}
	log.Println("Message service initialized successfully")

	// Fetch messages from the SQS queue
	queueUrl := os.Getenv("QueueUrl")

	for {

		messages, err := messageService.FetchMessages(queueUrl, 10, 20)
		if err != nil {
			log.Fatalf("Failed to fetch messages, %v", err)
		}
		if len(messages) == 0 {
			log.Println("No new messages to process")
			continue
		}

		// Process messages
		for _, message := range messages {
			log.Printf("Message ID: %s, Body: %s", *message.MessageId, *message.Body)
			err := messageService.DeleteMessage(queueUrl, message.ReceiptHandle)
			if err != nil {
				log.Printf("Failed to delete message: %s, %v", *message.MessageId, err)
			} else {
				log.Printf("Successfully deleted message: %s", *message.MessageId)
			}
		}

	}
}
