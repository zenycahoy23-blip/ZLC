FROM jenkins/jenkins:lts-jdk17

USER root

RUN apt-get update && \
    apt-get install -y docker.io sudo && \
    rm -rf /var/lib/apt/lists/*

# Add jenkins to docker group and allow sudo for chmod
RUN usermod -aG docker jenkins && \
    echo "jenkins ALL=(ALL) NOPASSWD: /bin/chmod" >> /etc/sudoers

# Install docker-compose
RUN curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose && \
    chmod +x /usr/local/bin/docker-compose

# Install buildx
RUN mkdir -p /usr/local/lib/docker/cli-plugins && \
    curl -L "https://github.com/docker/buildx/releases/download/v0.17.0/buildx-v0.17.0.linux-amd64" \
    -o /usr/local/lib/docker/cli-plugins/docker-buildx && \
    chmod +x /usr/local/lib/docker/cli-plugins/docker-buildx

USER jenkins