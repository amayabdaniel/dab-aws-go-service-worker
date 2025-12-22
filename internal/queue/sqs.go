package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/amayabdaniel/dab-aws-go-service-worker/pkg/config"
)

type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

type JobMessage struct {
	JobID string `json:"job_id"`
}

func NewSQSClient(cfg *config.Config) (*SQSClient, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var sqsClient *sqs.Client
	if cfg.SQSEndpoint != "" {
		sqsClient = sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(cfg.SQSEndpoint)
		})
	} else {
		sqsClient = sqs.NewFromConfig(awsCfg)
	}

	return &SQSClient{
		client:   sqsClient,
		queueURL: cfg.SQSQueueURL,
	}, nil
}

func (s *SQSClient) SendMessage(ctx context.Context, jobID string) error {
	message := JobMessage{JobID: jobID}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &s.queueURL,
		MessageBody: aws.String(string(body)),
	})
	return err
}

func (s *SQSClient) ReceiveMessages(ctx context.Context) ([]types.Message, error) {
	result, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, err
	}
	return result.Messages, nil
}

func (s *SQSClient) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.queueURL,
		ReceiptHandle: &receiptHandle,
	})
	return err
}

func (s *SQSClient) CreateQueueIfNotExists(ctx context.Context, queueName string) error {
	_, err := s.client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: &queueName,
	})
	return err
}