package service

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gopheramol/learn-aws/client"
)

type MessageService interface {
	FetchMessages(queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error)
	DeleteMessage(queueUrl string, receiptHandle *string) error
}

type messageService struct {
	awsClient client.AWSClient
}

func NewMessageService(awsClient client.AWSClient) MessageService {
	if awsClient == nil {
		log.Fatalf("AWSClient is not initialized")
	}
	return &messageService{
		awsClient: awsClient,
	}
}

func (s *messageService) FetchMessages(queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error) {
	if s.awsClient == nil {
		log.Fatalf("AWSClient is nil in MessageService")
	}
	return s.awsClient.GetMessages(queueUrl, maxMessages, waitTime)
}

func (s *messageService) DeleteMessage(queueUrl string, receiptHandle *string) error {
	if s.awsClient == nil {
		log.Fatalf("AWSClient is nil in MessageService")
	}
	return s.awsClient.DeleteMessage(queueUrl, receiptHandle)
}
