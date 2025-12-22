#!/bin/bash

echo "Starting services locally..."

# Work around permission issues
export DOCKER_HOST=unix:///var/run/docker.sock

# Try to start services
echo "Note: If you get permission errors, run this in a new terminal after logging out/in"
echo "Or run with: sudo ./run-local.sh"

docker-compose up -d

if [ $? -eq 0 ]; then
    echo "Services starting..."
    echo ""
    echo "Wait 10 seconds for services to start, then test with:"
    echo "  curl http://localhost:8080/health"
    echo ""
    echo "View logs with:"
    echo "  docker-compose logs -f"
else
    echo ""
    echo "Failed to start. Try:"
    echo "1. Make sure Docker Desktop is running"
    echo "2. Open a new terminal (the docker group needs a fresh session)"
    echo "3. Run: cd ~/projects/dab-aws-go-service-worker && docker-compose up -d"
fi