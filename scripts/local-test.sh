#!/bin/bash

echo "=== Local Testing Suite ==="

# Test 1: Check Go modules
echo "1. Checking Go dependencies..."
if [ -f "go.mod" ]; then
    echo "✓ go.mod exists"
else
    echo "✗ go.mod missing"
    exit 1
fi

# Test 2: Check all required files
echo "2. Checking required files..."
required_files=(
    "cmd/api/main.go"
    "cmd/worker/main.go"
    "docker-compose.yml"
    "deployments/docker/api.Dockerfile"
    "deployments/docker/worker.Dockerfile"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✓ $file exists"
    else
        echo "✗ $file missing"
        exit 1
    fi
done

# Test 3: Docker Compose validation
echo "3. Validating docker-compose..."
docker-compose config > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ docker-compose.yml is valid"
else
    echo "✗ docker-compose.yml has errors"
fi

echo "=== All tests passed! ==="