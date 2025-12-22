#!/bin/bash

echo "=== Initializing Local Environment ==="

# Create SQS queue in LocalStack
echo "Creating SQS queue..."
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name jobs-queue --region us-east-2

# List queues to verify
echo "Verifying queue creation..."
aws --endpoint-url=http://localhost:4566 sqs list-queues --region us-east-2

echo "=== Local environment ready ==="