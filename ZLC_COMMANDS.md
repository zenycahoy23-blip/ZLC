# ZLC System - Quick Reference Commands

## Starting & Stopping

# Start all services
docker compose up

# Start in background
docker compose up -d

# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v

# Rebuild all images
docker compose build --no-cache

## Viewing Logs

# View all logs
docker compose logs

# Follow logs in real-time
docker compose logs -f

# View specific service logs
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres

## Container Management

# List running containers
docker compose ps

# Execute command in container
docker compose exec backend bash
docker compose exec frontend bash
docker compose exec postgres psql -U authuser -d authdb

# View container info
docker inspect zlc-backend
docker inspect zlc-frontend

## Database Access

# Connect to PostgreSQL
docker compose exec postgres psql -U authuser -d authdb

# Useful psql commands:
# \dt              - list tables
# \d users         - describe users table
# SELECT * FROM users;
# SELECT * FROM audit_logs;
# SELECT * FROM login_history;

## Metrics & Monitoring

# Prometheus metrics (raw)
curl http://localhost:9092/metrics

# View Prometheus targets
curl http://localhost:9092/api/v1/targets

## Network & Ports

# Check port usage (Linux)
netstat -tulpn | grep :9031
netstat -tulpn | grep :8084
netstat -tulpn | grep :3002

# Check network connections
docker network inspect zlc-network
docker inspect zlc-backend | grep Networks

## System Health

# View resource usage
docker stats

# View disk usage
docker system df

# Clean up unused resources
docker system prune -a

## Troubleshooting

# Restart a service
docker compose restart backend
docker compose restart frontend

# Remove a container and recreate
docker compose rm -f backend
docker compose up -d backend

# View build output
docker compose build --no-cache 2>&1 | tail -100

# Check service dependencies
docker compose ps --services --filter "status=running"

## Development

# Rebuild frontend only
docker compose build frontend

# Rebuild backend only
docker compose build backend

# Watch frontend logs
docker compose logs -f frontend | grep -i error

# Monitor API calls
curl -v http://localhost:8084/api/health

## Useful Links

Frontend:      http://localhost:9031
Backend API:   http://localhost:8084
Grafana:       http://localhost:3002
Jenkins:       http://localhost:8082
Prometheus:    http://localhost:9092
Database:      localhost:5434

Default Admin:
Email:    cliffe026@gmail.com
Password: admin123

SMTP Config:
Email: zenycahoy23@gmail.com
App:   ZLC
Pass:  jquq chgi wzsy jjfq
