#!/bin/bash

echo "=== Verifying Project Structure ==="

# Check imports in main files
echo "Checking API imports..."
grep -E "github.com/amayabdaniel/dab-aws-go-service-worker" cmd/api/main.go > /dev/null
if [ $? -eq 0 ]; then
    echo "✓ API imports look correct"
else
    echo "✗ API import issues"
fi

echo "Checking Worker imports..."
grep -E "github.com/amayabdaniel/dab-aws-go-service-worker" cmd/worker/main.go > /dev/null
if [ $? -eq 0 ]; then
    echo "✓ Worker imports look correct"
else
    echo "✗ Worker import issues"
fi

# Check package declarations
echo "Checking package structure..."
for pkg in config logger database queue models; do
    if [ -d "internal/$pkg" ] || [ -d "pkg/$pkg" ]; then
        echo "✓ Package $pkg exists"
    else
        echo "✗ Package $pkg missing"
    fi
done

echo "=== Structure verification complete ==="