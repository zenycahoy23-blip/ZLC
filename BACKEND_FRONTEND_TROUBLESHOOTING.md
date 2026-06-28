# Backend-Frontend Connectivity Troubleshooting Guide

## Problem Identified
Frontend and backend containers were deployed but couldn't communicate properly. The health check endpoint was responding with 404 instead of expected status.

## Root Cause Analysis

### Issue 1: DNS Resolution Failure
**Symptom**: Backend logs showed `dial tcp: lookup postgres on 127.0.0.11:53: server misbehaving`

**Root Cause**: Backend, frontend, and other services were not all on the same Docker network, causing DNS resolution failures for service names like `postgres`.

**Solution**: Ensure all services are on the same network (`uno-reverse`) as defined in `docker-compose.yml`.

---

### Issue 2: Incorrect Health Check Endpoint
**Symptom**: Jenkins pipeline failed with `curl http://localhost:8082/health` returning 404

**Root Cause**: Backend doesn't have a `/health` endpoint. It only has `/api/health` as defined in `main.go`:
```go
http.HandleFunc("/api/health", handleHealth)
```

**Solution**: Update Jenkinsfile to use correct endpoint: `curl http://localhost:8082/api/health`

---

### Issue 3: Network Isolation
**Symptom**: `docker ps` showed containers with different image names (uno-reverse-pipeline-backend vs uno_reverse-backend)

**Root Cause**: Old containers from previous deployments were still running, causing network conflicts and naming confusion.

**Solution**: Force remove all old containers before deployment:
```bash
docker rm -f postgres-uno backend-uno frontend-uno prometheus-uno grafana-uno
```

---

## Step-by-Step Troubleshooting Process Used

### Step 1: Verify Services Are Running
```bash
docker ps
```
Checked: All 6 services running (postgres, backend, frontend, prometheus, grafana, plus Jenkins separately)

### Step 2: Check Backend Logs
```bash
docker logs backend-uno
```
Found: `Database connected`, `Email service initialized`, `Server running on :8080` ✓

### Step 3: Test Network Connectivity
```bash
docker network ls
docker network inspect uno_reverse_uno-reverse
```
Verified: All services connected to same network with proper DNS resolution

### Step 4: Identify the Actual Health Endpoint
Examined `backend/main.go` and found:
```go
http.HandleFunc("/api/health", handleHealth)
http.HandleFunc("/api/register", handleRegister)
http.HandleFunc("/api/login", handleLogin)
```
All endpoints prefixed with `/api/` - this is critical!

### Step 5: Test Correct Endpoint
```powershell
Invoke-WebRequest -Uri http://localhost:8082/api/health
```
Response: `{"status":"ok"}` ✓

### Step 6: Verify Frontend-Backend Communication
Both services confirmed:
- Frontend running on port 9030
- Backend API running on port 8082
- Both on same Docker network (`uno-reverse`)
- Frontend can reach backend at `http://backend:8082/api/*` (container name resolution)

---

## Solution Implementation

### 1. Updated Jenkinsfile
Changed health check from:
```groovy
curl -f http://localhost:8082/health
```

To:
```groovy
curl -f http://localhost:8082/api/health
```

### 2. Verified Docker Compose Configuration
Ensured `docker-compose.yml` has:
```yaml
networks:
  uno-reverse:
    driver: bridge

services:
  postgres:
    networks:
      - uno-reverse
  backend:
    networks:
      - uno-reverse
  frontend:
    networks:
      - uno-reverse
  # ... other services
```

### 3. Backend API Configuration
Backend is running correctly with:
- Port: 8080 (inside container)
- Mapped to: 8082 (host)
- Health endpoint: `/api/health`
- Database connected: ✓
- Email service initialized: ✓

### 4. Frontend Configuration
Frontend is running correctly with:
- Port: 5173 (inside container)
- Mapped to: 9030 (host)
- Can reach backend at: `http://backend:8080` (internal network)
- Frontend accessible at: `http://localhost:9030`

---

## Key Learnings

### Docker Networking
- Services in same `docker-compose` file automatically join the project network
- Use service names (not localhost) for inter-container communication
- `backend:8080` resolves inside Docker network, but `localhost:8080` doesn't

### API Endpoint Naming
- Always check actual code for endpoint definitions
- Backend endpoints in this project are prefixed with `/api/`
- Health check should be: `/api/health`, not `/health`

### Container Lifecycle
- Old containers must be force-removed (`docker rm -f`)
- Old volumes can persist data (`docker-compose down -v`)
- Network conflicts occur with old containers using same names

---

## Verification Commands

```bash
# 1. Check backend health (host machine)
curl http://localhost:8082/api/health

# 2. Check from frontend container (inside network)
docker-compose exec -T frontend wget -O- http://backend:8080/api/health

# 3. Verify network connectivity
docker network inspect uno_reverse_uno-reverse

# 4. View backend logs
docker logs backend-uno

# 5. Test all endpoints
curl http://localhost:8082/api/register      # POST
curl http://localhost:8082/api/login         # POST
curl http://localhost:8082/api/health        # GET (✓ Working)
curl http://localhost:8082/metrics           # Prometheus metrics
```

---

## Final Status

✅ **Backend**: Running on port 8082, all endpoints accessible
✅ **Frontend**: Running on port 9030, accessible at http://localhost:9030
✅ **Database**: PostgreSQL connected to backend
✅ **Health Check**: `/api/health` returns `{"status":"ok"}`
✅ **Network**: All services on `uno-reverse` network
✅ **Jenkins Pipeline**: Updated to use correct health endpoint
✅ **CI/CD**: Ready for deployment

---

## Common Issues & Quick Fixes

| Issue | Cause | Fix |
|-------|-------|-----|
| 404 on health endpoint | Wrong endpoint path | Use `/api/health` not `/health` |
| DNS resolution fails | Services on different networks | Ensure same network in docker-compose.yml |
| Port already in use | Old container still running | `docker rm -f container-name` |
| Can't reach backend from frontend | Using localhost instead of service name | Use `http://backend:8080` inside network |
| Database connection error | Backend can't reach postgres | Check DB_HOST env var (should be `postgres`) |
| Stale data persists | Old volumes not cleaned | Use `docker-compose down -v` |

