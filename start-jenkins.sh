#!/bin/bash
# Standalone Jenkins Setup for UnoReverse Pipeline
# Run this script to start Jenkins separately from docker-compose

echo "Starting standalone Jenkins for UnoReverse..."

# Create jenkins_home volume if it doesn't exist
docker volume create jenkins_home 2>/dev/null || true

# Stop and remove old jenkins-uno if it exists
docker rm -f jenkins-uno 2>/dev/null || true
sleep 2

# Build Jenkins image from Dockerfile
docker build -t uno_reverse-jenkins ./jenkins

# Run Jenkins on the uno-reverse network
docker run -d \
  --name jenkins-uno \
  --network uno-reverse \
  --privileged \
  --user root \
  -p 8081:8080 \
  -p 50000:50000 \
  -e JENKINS_OPTS="-Djenkins.install.runSetupWizard=false" \
  -v jenkins_home:/var/jenkins_home \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --restart unless-stopped \
  --health-cmd="curl -f http://localhost:8080/cli/ || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=5 \
  --health-start-period=60s \
  uno_reverse-jenkins

echo "Jenkins started successfully!"
echo "Access Jenkins at: http://localhost:8081"
echo ""
echo "To view Jenkins logs: docker logs jenkins-uno"
echo "To stop Jenkins: docker stop jenkins-uno"
echo "To restart Jenkins: docker restart jenkins-uno"
