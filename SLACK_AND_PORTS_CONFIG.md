# UnoReverse Configuration Summary

## 🔄 Port Changes

All backend services have been updated to avoid port collisions:

| Service | Old Port | New Port | Access |
|---------|----------|----------|--------|
| PostgreSQL | 5432 | **5433** | `localhost:5433` |
| Backend API | 8082 | **8083** | `http://localhost:8083` |
| Prometheus | 9090 | **9091** | `http://localhost:9091` |
| Grafana | 3000 | **3001** | `http://localhost:3001` (admin/admin123) |
| Frontend | 9030 | 9030 | `http://localhost:9030` |
| Jenkins | 8081 | 8081 | `http://localhost:8081` |

## 📢 Slack Integration

The Jenkinsfile has been updated with Slack notifications at every stage:

- 🚀 **Build Start**: Notifies when build begins
- 🔨 **Build Stage**: Notifies when Docker images start building
- 📦 **Deploy Stage**: Notifies when deployment begins
- 🏥 **Health Check**: Notifies when health checks start
- ✅ **Success**: Full URLs to all services on successful build
- ❌ **Failure**: Link to Jenkins console on failure
- ⚠️ **Unstable**: Alert for unstable builds

### Setting Up Slack Webhook

1. **Create a Slack App:**
   - Go to https://api.slack.com/apps
   - Click "Create New App"
   - Select "From scratch"
   - Name your app (e.g., "UnoReverse Jenkins")
   - Select your workspace

2. **Enable Incoming Webhooks:**
   - Navigate to "Incoming Webhooks" in the left sidebar
   - Toggle "Activate Incoming Webhooks" to ON
   - Click "Add New Webhook to Workspace"
   - Select the channel to post to (e.g., #deployments)
   - Authorize

3. **Copy Your Webhook URL:**
   - You'll get a URL like: `https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX`

4. **Add to Jenkins:**
   - Go to `http://localhost:8081/credentials/store/system/domain/_/newSecret`
   - Select "Secret text"
   - Paste the webhook URL in the "Secret" field
   - Set ID to: `slack-webhook`
   - Description: "Slack Webhook for Notifications"
   - Save

### Install Slack Plugin

1. Go to Jenkins: `Manage Jenkins` → `Plugins` → `Available plugins`
2. Search for "Slack"
3. Install the Slack plugin by jbaruch
4. Restart Jenkins

## 🚀 Updated Services

All services in `docker-compose.yml` now use the new ports. Key changes:

- **Frontend** now connects to backend on port 8083
- **Health checks** reference new port numbers
- **Jenkinsfile** updated with new port references

## 📝 Files Modified

1. **docker-compose.yml**
   - PostgreSQL: 5433
   - Backend: 8083
   - Prometheus: 9091
   - Grafana: 3001

2. **Jenkinsfile**
   - Added Slack webhook environment variable
   - Added notifications to all stages
   - Updated health check ports
   - Enhanced post-build notifications

3. **setup-slack.sh**
   - New setup script for Slack configuration
   - Instructions for webhook setup

## ✅ Verification

After making these changes:

```bash
# Verify port mappings
docker-compose ps

# Test the backend
curl http://localhost:8083/api/health

# View Jenkins
http://localhost:8081

# View Grafana
http://localhost:3001

# View Prometheus
http://localhost:9091
```

## 🔧 Troubleshooting

**Slack notifications not sending?**
- Verify the webhook URL is correct in Jenkins credentials
- Check Jenkins logs: `docker logs jenkins-uno`
- Ensure Slack plugin is installed

**Backend not responding on 8083?**
- Check if container is running: `docker-compose ps`
- View logs: `docker-compose logs backend`
- Verify environment variables are correct

**Port still in use?**
- Check what's using the port: `netstat -tulpn | grep :<port>`
- Kill the process or restart Docker
