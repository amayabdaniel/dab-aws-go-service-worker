#!/bin/bash

echo "Testing Go compilation..."

# Test API compilation
echo "Compiling API service..."
go build -o /tmp/api-test cmd/api/main.go
if [ $? -eq 0 ]; then
    echo "✅ API compiled successfully"
    rm /tmp/api-test
else
    echo "❌ API compilation failed"
    exit 1
fi

# Test Worker compilation
echo "Compiling Worker service..."
go build -o /tmp/worker-test cmd/worker/main.go
if [ $? -eq 0 ]; then
    echo "✅ Worker compiled successfully"
    rm /tmp/worker-test
else
    echo "❌ Worker compilation failed"
    exit 1
fi

echo "All services compile successfully!"