#!/bin/bash

# Jenkins Setup Script for Uno Reverse Authentication System
# This script configures Jenkins with all necessary plugins and credentials

set -e

JENKINS_HOME=${JENKINS_HOME:-/var/jenkins_home}
JENKINS_CLI=${JENKINS_HOME}/jenkins-cli.jar
JENKINS_URL=${JENKINS_URL:-http://localhost:8080}

echo "=== Jenkins Setup for Uno Reverse Authentication System ==="

# Download Jenkins CLI
if [ ! -f ${JENKINS_CLI} ]; then
    echo "Downloading Jenkins CLI..."
    wget -q ${JENKINS_URL}/jnlpJars/jenkins-cli.jar -O ${JENKINS_CLI}
fi

# Wait for Jenkins to be ready
echo "Waiting for Jenkins to be ready..."
until curl -s ${JENKINS_URL}/api/json > /dev/null; do
    echo "Jenkins is starting up..."
    sleep 10
done

echo "Jenkins is ready!"

# Install required plugins
echo "Installing required plugins..."
java -jar ${JENKINS_CLI} -s ${JENKINS_URL} install-plugin \
    docker-plugin \
    docker-build-step \
    docker-commons \
    git \
    github \
    github-branch-source \
    email-ext \
    pipeline-model-definition \
    pipeline-stage-view \
    timestamper \
    log-parser \
    ansicolor \
    warnings-ng \
    performance \
    cobertura \
    htmlpublisher \
    junit \
    -restart

echo "Waiting for plugins to be installed and Jenkins to restart..."
sleep 60

# Wait for Jenkins to restart
until curl -s ${JENKINS_URL}/api/json > /dev/null; do
    echo "Waiting for Jenkins to restart..."
    sleep 10
done

echo "=== Jenkins Setup Complete ==="
echo ""
echo "Next steps:"
echo "1. Access Jenkins at ${JENKINS_URL}"
echo "2. Add credentials in Jenkins UI:"
echo "   - GitHub credentials (for repository access)"
echo "   - Docker Hub credentials (for image registry)"
echo "3. Create a new Pipeline job"
echo "4. Point it to your repository with Jenkinsfile"
echo "5. Configure webhook on GitHub to trigger builds"
