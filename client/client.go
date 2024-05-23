package client

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// AWSClient is the interface for interacting with AWS SQS.
type AWSClient interface {
	GetMessages(queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error)
	DeleteMessage(queueUrl string, receiptHandle *string) error
}

// awsClient implements the AWSClient interface.
type awsClient struct {
	awsConfig aws.Config
	sqsClient *sqs.Client
}

// NewAWSClient creates a new instance of awsClient and returns it as an AWSClient.
func NewAWSClient(awsConfig aws.Config) AWSClient {
	return &awsClient{
		awsConfig: awsConfig,
		sqsClient: sqs.NewFromConfig(awsConfig),
	}
}

// GetMessages fetches messages from the specified SQS queue.
func (client *awsClient) GetMessages(queueUrl string, maxMessages int32, waitTime int32) ([]types.Message, error) {
	if client.sqsClient == nil {
		log.Fatalf("SQS client is not initialized")
	}
	if queueUrl == "" {
		log.Fatalf("Queue URL is empty")
	}

	result, err := client.sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})
	if err != nil {
		log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", queueUrl, err)
		return nil, err
	}
	return result.Messages, nil
}

func (client *awsClient) DeleteMessage(queueUrl string, receiptHandle *string) error {
	_, err := client.sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: receiptHandle,
	})
	return err
}
