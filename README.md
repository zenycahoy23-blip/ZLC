# ZLC - Secure Authentication System

ZLC is an enterprise-grade authentication system with advanced security features, monitoring, and CI/CD integration.

## Features

- **Secure Authentication**: Email & phone verification with OTP
- **Role-Based Access Control**: Admin, Moderator, and User roles
- **Audit Logging**: Comprehensive activity tracking with device & location info
- **Monitoring**: Real-time metrics with Prometheus and Grafana
- **CI/CD**: Jenkins integration for automated deployments
- **Multi-factor Security**: Account locking, password history, and breach prevention

## Services

- **Backend**: Go REST API with PostgreSQL
- **Frontend**: Svelte UI with responsive design
- **Database**: PostgreSQL for data persistence
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: Jenkins

## Quick Start

```bash
docker compose up --build
```

### Access Points

- **Frontend**: http://localhost:9031
- **Backend API**: http://localhost:8084
- **Grafana**: http://localhost:3002 (admin/admin123)
- **Jenkins**: http://localhost:8082
- **Prometheus**: http://localhost:9092

### Database Credentials

- User: authuser
- Password: authpass
- Database: authdb
- Port: 5434

## Default Admin Account

- Email: cliffe026@gmail.com
- Password: admin123
- Note: Change these credentials in production

## Configuration

Edit `.env` file to customize:
- Database settings
- JWT secret
- SMTP credentials (Gmail configured)
- OTP expiry time

## Architecture

```
Frontend (Svelte) → Backend (Go) → PostgreSQL
                 ↓
            Prometheus/Grafana
                 ↓
              Jenkins
```

## Security Features

- Password encryption with bcrypt
- JWT token-based authentication
- OTP verification for phone numbers
- Account lockout on multiple failed attempts
- Detailed audit logging
- IP and device tracking

## License

Proprietary - ZLC Systems
