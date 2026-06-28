#!/bin/bash

# Jenkins Slack Configuration Script
# This script sets up Slack integration in Jenkins

echo "📢 Setting up Slack Integration for Jenkins..."

# Jenkins URL and credentials
JENKINS_URL="http://localhost:8081"
JENKINS_USER="admin"
JENKINS_TOKEN="${JENKINS_API_TOKEN:-your-jenkins-token}"

# Slack webhook URL (you need to set this in Jenkins)
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

echo ""
echo "======================================"
echo "📌 SLACK WEBHOOK SETUP"
echo "======================================"
echo ""
echo "1. Go to your Slack workspace at https://api.slack.com/apps"
echo "2. Create a new app or select existing one"
echo "3. Enable 'Incoming Webhooks'"
echo "4. Click 'Add New Webhook to Workspace'"
echo "5. Select a channel and authorize"
echo "6. Copy the Webhook URL"
echo ""
echo "======================================"
echo ""
echo "To add the credential to Jenkins manually:"
echo "1. Go to: http://localhost:8081/credentials/store/system/domain/_/newSecret"
echo "2. Select: Secret text"
echo "3. Paste your Slack Webhook URL in 'Secret' field"
echo "4. Set ID to: slack-webhook"
echo "5. Save"
echo ""
echo "======================================"
echo ""

# Try to add credential via Jenkins API (requires authentication)
if [ ! -z "$JENKINS_TOKEN" ] && [ "$JENKINS_TOKEN" != "your-jenkins-token" ]; then
    echo "Attempting to configure Slack webhook via API..."
    
    curl -u "${JENKINS_USER}:${JENKINS_TOKEN}" -X POST \
    "${JENKINS_URL}/credentials/store/system/domain/_/createCredentials" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "json={\"\":\"0\",\"credentials\":{\"scope\":\"GLOBAL\",\"id\":\"slack-webhook\",\"secret\":\"${SLACK_WEBHOOK_URL}\",\"description\":\"Slack Webhook for Notifications\",\"\$class\":\"org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl\"}}" \
    && echo "✅ Slack webhook credential added successfully!" || echo "⚠️ Could not auto-add credential. Please add manually via Jenkins UI."
else
    echo "⚠️ JENKINS_API_TOKEN not set. Please configure manually via Jenkins UI."
fi

echo ""
echo "======================================"
echo "📋 NEXT STEPS:"
echo "======================================"
echo ""
echo "1. Ensure Slack plugin is installed in Jenkins:"
echo "   - Go to: Manage Jenkins > Plugin Manager"
echo "   - Search for 'Slack' and install"
echo ""
echo "2. The Jenkinsfile now includes:"
echo "   - Build start notification"
echo "   - Build stage notifications (🔨 Build, 📦 Deploy, 🏥 Health Check)"
echo "   - Success/Failure/Unstable notifications"
echo "   - Links to Jenkins console and services"
echo ""
echo "3. Your backend ports have been changed:"
echo "   - PostgreSQL: 5433 (was 5432)"
echo "   - Backend API: 8083 (was 8082)"
echo "   - Prometheus: 9091 (was 9090)"
echo "   - Grafana: 3001 (was 3000)"
echo "   - Jenkins: 8081 (unchanged)"
echo ""
echo "======================================"
echo ""
