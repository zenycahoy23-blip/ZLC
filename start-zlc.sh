#!/bin/bash

# ZLC System - Docker Compose Start Script
# This script starts all ZLC microservices

echo "Starting ZLC Authentication System..."
echo "======================================"
echo ""

# Check Docker daemon
if ! docker info > /dev/null 2>&1; then
  echo "❌ Docker daemon is not running. Please start Docker Desktop or daemon."
  exit 1
fi

echo "✓ Docker is running"

# Pull latest images if needed
echo "Building ZLC services..."
docker compose build --pull

if [ $? -ne 0 ]; then
  echo "❌ Build failed. Check logs above."
  exit 1
fi

echo "✓ Build successful"
echo ""
echo "Starting containers..."
docker compose up -d

if [ $? -ne 0 ]; then
  echo "❌ Failed to start containers."
  exit 1
fi

echo "✓ All containers started"
echo ""
echo "======================================"
echo "ZLC System is running!"
echo "======================================"
echo ""
echo "📱 Access Points:"
echo "  Frontend:    http://localhost:9031"
echo "  Backend API: http://localhost:8084"
echo "  Grafana:     http://localhost:3002 (admin/admin123)"
echo "  Jenkins:     http://localhost:8082"
echo "  Prometheus:  http://localhost:9092"
echo "  Database:    localhost:5434"
echo ""
echo "🔑 Default Credentials:"
echo "  Email:    cliffe026@gmail.com"
echo "  Password: admin123"
echo ""
echo "⚠️  Change these credentials in production!"
echo ""
echo "📊 Monitor logs:"
echo "  docker compose logs -f [service_name]"
echo ""
echo "Stop system:"
echo "  docker compose down"
echo ""
