# UnoReverse Jenkins + Docker Compose Best Practices

## Why the Issues Happen

### 1. **Old Containers Not Removed**
- `docker-compose up -d` will NOT recreate containers if they already exist
- Old image layers are cached in Docker
- Solution: Use `--force-recreate` and `docker-compose down` first

### 2. **Docker Image Layer Caching**
- Docker caches layers from previous builds
- `docker-compose build` reuses cached layers if Dockerfile hasn't changed
- Old code in frontend/backend still runs from cache
- Solution: Use `docker-compose build --no-cache`

### 3. **Stale Frontend Code**
- Browser caches old JavaScript/CSS files
- Frontend image may not update properly
- Solution: Use `--no-cache` builds + browser hard refresh (Ctrl+Shift+R)

### 4. **Git Not Pulling Latest Code**
- Jenkins workspace may have old code
- `checkout scm` might use cached Git data
- Solution: Use `deleteDir()` before checkout + `clean: true, cleanBeforeCheckout: true`

### 5. **Container Naming Conflicts**
- Can't start new container if old one still exists
- Old volumes still mounted
- Solution: `docker-compose down --remove-orphans` before `up`

### 6. **Jenkins Container Conflict (Fixed)**
- If Jenkins was in docker-compose, `docker-compose down` would stop Jenkins
- Pipeline would fail mid-execution
- Solution: Run Jenkins standalone on the network

### 7. **Dangling Volumes**
- Old database volumes persist with stale data
- New containers connect to old volumes
- Solution: Use `-v` flag in `docker-compose down` or `docker volume prune`

---

## The Complete Safe Deployment Flow

```
1. WORKSPACE CLEANUP
   └─ deleteDir() - removes all old Jenkins files

2. CODE CHECKOUT
   └─ Fresh git clone with clean: true

3. PRE-DEPLOY CLEANUP
   ├─ docker-compose down --remove-orphans
   ├─ docker rm -f (force remove lingering containers)
   └─ docker volume prune -f (remove old volumes)

4. BUILD FRESH
   ├─ docker-compose build --no-cache (no layer caching)
   └─ docker image prune -f (cleanup old images)

5. DEPLOY FRESH
   ├─ docker-compose up -d --force-recreate (recreate all)
   └─ wait for containers to stabilize

6. HEALTH CHECKS
   ├─ PostgreSQL ready
   ├─ Backend API responding
   └─ Frontend accessible

7. SUCCESS
   └─ Fresh code running with no conflicts
```

---

## Jenkinsfile Key Features

### `deleteDir()`
- Cleans entire Jenkins workspace
- Ensures no stale files from previous builds
- Forces fresh checkout

### `clean: true, cleanBeforeCheckout: true`
- Git option to clean untracked files
- Prevents old code from persisting

### `disableConcurrentBuilds()`
- Prevents two deployments running simultaneously
- Avoids port conflicts and race conditions

### `--no-cache`
- Forces Docker to rebuild from scratch
- No layer caching
- Guarantees latest code

### `--force-recreate`
- Removes and recreates every container
- Even if image hasn't changed
- Ensures fresh environment

### `--remove-orphans`
- Removes containers defined in compose but not in current config
- Prevents lingering old services

### `docker volume prune -f`
- Removes all unused volumes
- Prevents old database data

---

## Docker Compose Best Practices

### 1. **Always Use Named Containers**
```yaml
services:
  backend:
    container_name: backend-uno  # Explicit name
    ports:
      - "8082:8080"
```

### 2. **Use Explicit Networks**
```yaml
networks:
  uno-reverse:
    driver: bridge
```

### 3. **Health Checks**
```yaml
postgres:
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U authuser -d authdb"]
    interval: 10s
    timeout: 5s
    retries: 5
```

### 4. **Volumes Strategy**
```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data  # Named volume
  - ./backend/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql  # Bind mount
```

### 5. **Environment Variables**
```yaml
environment:
  DB_HOST: postgres              # Service name, not localhost
  DB_USER: ${DB_USER:-authuser}  # Use defaults
```

---

## Frontend Caching Issues & Solutions

### Problem
Browser caches old JavaScript/CSS files even after deployment

### Solutions

#### 1. **Hard Browser Refresh**
- Chrome/Edge: `Ctrl+Shift+R` (Windows) or `Cmd+Shift+R` (Mac)
- Firefox: `Ctrl+Shift+R` (Windows) or `Cmd+Shift+R` (Mac)
- Safari: `Cmd+Shift+R`

#### 2. **Clear Browser Cache**
- Chrome DevTools: F12 → Application → Clear storage

#### 3. **Cache Busting in Frontend**
Add to `vite.config.ts` or `webpack.config.js`:
```javascript
build: {
  rollupOptions: {
    output: {
      entryFileNames: `[name].[hash].js`,
      chunkFileNames: `[name].[hash].js`,
      assetFileNames: `[name].[hash].[ext]`
    }
  }
}
```

#### 4. **Add Version Header**
In backend response headers:
```
Cache-Control: no-cache, no-store, must-revalidate
Pragma: no-cache
Expires: 0
```

---

## Manual Cleanup Commands

```bash
# Full system cleanup
./cleanup.sh

# Or manually:

# Stop all services
docker-compose down --remove-orphans

# Remove specific containers
docker rm -f postgres-uno backend-uno frontend-uno prometheus-uno grafana-uno

# Remove all dangling volumes
docker volume prune -f

# Remove all unused images
docker image prune -f

# See disk space
docker system df

# Deep clean (caution: affects all Docker)
docker system prune -a -f
```

---

## Troubleshooting

### Issue: "Container already in use"
```bash
docker rm -f container-name
```

### Issue: "Port already in use"
```bash
docker ps -a
docker stop container-id
docker rm container-id
```

### Issue: Old data persists
```bash
docker volume rm uno_reverse_postgres_data
docker-compose down -v  # -v removes volumes
```

### Issue: Stale frontend code
```bash
# Hard refresh browser (Ctrl+Shift+R)
# Or clear cache manually
# Rebuild without cache
docker-compose build --no-cache
```

### Issue: Image layer caching preventing updates
```bash
docker-compose build --no-cache --force-rm
```

---

## Production Checklist

- [x] Jenkins runs separately (not in docker-compose)
- [x] `deleteDir()` clears workspace
- [x] `--no-cache` prevents stale layers
- [x] `--force-recreate` ensures fresh containers
- [x] Health checks validate deployments
- [x] Volumes cleaned with `-v` flag
- [x] Network explicitly defined
- [x] Container names explicit
- [x] Environment variables use defaults
- [x] disableConcurrentBuilds() prevents race conditions
- [x] Post-failure logs capture all errors
- [x] Frontend cache-busting implemented
- [x] Database migrations run on startup

---

## Jenkins Pipeline Execution Order

```
1. Cleanup Workspace          (removes old Jenkins files)
2. Checkout Fresh Code        (git clone with clean)
3. Pre-Deploy Cleanup         (docker-compose down, prune)
4. Build Fresh Images         (--no-cache, --force-rm)
5. Deploy Fresh Containers    (--force-recreate)
6. Health Check - Database    (PostgreSQL ready)
7. Health Check - Backend     (API responds)
8. Health Check - Frontend    (HTTP 200)
9. Verify Services            (final status)
10. Post Actions              (success/failure logs)
```

---

## How to Use

### First Deployment
```bash
cd /path/to/UnoReverse
docker-compose up -d
```

### After Pushing Code Changes
Jenkins automatically triggers:
```bash
# Jenkins pulls new code
git pull

# Workspace and containers cleaned
# Fresh images built
# Fresh containers deployed

# Result: Latest code running
```

### Manual Deployment
```bash
./cleanup.sh
docker-compose build --no-cache
docker-compose up -d
```

### Check Status
```bash
docker-compose ps
docker-compose logs backend
```
