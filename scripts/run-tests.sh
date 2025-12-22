#!/bin/bash

echo "=== RUNNING GO TESTS ==="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test results
total_tests=0
passed_tests=0
failed_tests=0

echo -e "\n${YELLOW}Running unit tests...${NC}"

# Run tests with coverage
go test -v -cover -coverprofile=coverage.out ./... 2>&1 | while IFS= read -r line; do
    echo "$line"
    
    # Count test results
    if [[ $line == *"PASS"* ]]; then
        ((passed_tests++))
    elif [[ $line == *"FAIL"* ]] && [[ $line != *"FAIL"*"github.com"* ]]; then
        ((failed_tests++))
    fi
    
    if [[ $line == *"RUN"* ]]; then
        ((total_tests++))
    fi
done

# Check if tests ran at all
if [ $? -ne 0 ]; then
    echo -e "\n${RED}Tests failed to run. Check Go installation and dependencies.${NC}"
    echo "Try running: go mod download"
    exit 1
fi

# Generate coverage report
if [ -f coverage.out ]; then
    echo -e "\n${YELLOW}Coverage Summary:${NC}"
    go tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $3}'
fi

# Run go vet
echo -e "\n${YELLOW}Running go vet...${NC}"
go vet ./...
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ No issues found by go vet${NC}"
else
    echo -e "${RED}✗ go vet found issues${NC}"
fi

# Run gofmt
echo -e "\n${YELLOW}Checking code formatting...${NC}"
unformatted=$(gofmt -l .)
if [ -z "$unformatted" ]; then
    echo -e "${GREEN}✓ All files properly formatted${NC}"
else
    echo -e "${RED}✗ Unformatted files:${NC}"
    echo "$unformatted"
fi

echo -e "\n${YELLOW}=== TEST SUMMARY ===${NC}"
echo "Tests will show detailed results above"
echo "Check for any compilation errors or test failures"

# Clean up
rm -f coverage.out