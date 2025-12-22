#!/bin/bash

# API Usage Examples

API_URL=${API_URL:-http://localhost:8080}

echo "=== API Health Check ==="
curl -s $API_URL/health | jq .

echo -e "\n=== Create a Data Processing Job ==="
JOB_ID=$(curl -s -X POST $API_URL/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "data-processing",
    "data": "Process this important business data with analysis"
  }' | jq -r .id)

echo "Created job: $JOB_ID"

echo -e "\n=== Wait for Processing ==="
sleep 3

echo -e "\n=== Check Job Status ==="
curl -s $API_URL/jobs/$JOB_ID | jq .

echo -e "\n=== Create a Health Report ==="
curl -s -X POST $API_URL/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "health-report",
    "data": "Generate system health metrics"
  }' | jq .

echo -e "\n=== Create a Batch Import Job ==="
curl -s -X POST $API_URL/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "batch-import",
    "data": "user1,john@example.com,John Doe\nuser2,jane@example.com,Jane Doe"
  }' | jq .

echo -e "\n=== List All Jobs ==="
curl -s $API_URL/jobs | jq .

echo -e "\n=== List Only Completed Jobs ==="
curl -s "$API_URL/jobs?status=completed" | jq .

echo -e "\n=== Invalid Request (Empty Type) ==="
curl -s -X POST $API_URL/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "",
    "data": "This should fail validation"
  }' | jq .