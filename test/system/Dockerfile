# System test environment for git-migrator
# This Dockerfile creates an environment with CVS and Git for end-to-end testing

FROM ubuntu:22.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
RUN apt-get update && apt-get install -y \
    cvs \
    git \
    golang-go \
    make \
    curl \
    vim \
    && rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
ENV GO111MODULE=on

# Create working directory
WORKDIR /app

# Copy the git-migrator source code
COPY . /app/

# Build the git-migrator tool
RUN make build

# Create test user (CVS requires a real user)
RUN useradd -m -s /bin/bash testuser && \
    echo "testuser:testuser" | chpasswd

# Set up CVS repository directory
ENV CVSROOT=/cvs-repo
RUN mkdir -p $CVSROOT && \
    cvs -d $CVSROOT init && \
    chown -R testuser:testuser $CVSROOT

# Create workspace for tests
RUN mkdir -p /workspace && chown -R testuser:testuser /workspace

# Set up environment for testuser
USER testuser
WORKDIR /home/testuser

# Configure git for testuser
RUN git config --global user.email "testuser@example.com" && \
    git config --global user.name "Test User" && \
    git config --global init.defaultBranch main

# Set up CVS environment variables for testuser
ENV CVSROOT=/cvs-repo
ENV CVSEDITOR=/usr/bin/vim

WORKDIR /workspace

# Copy test scripts
COPY test/system/setup-cvs-repo.sh /usr/local/bin/
COPY test/system/run-migration-test.sh /usr/local/bin/
COPY test/system/verify-migration.sh /usr/local/bin/

USER root
RUN chmod +x /usr/local/bin/*.sh
USER testuser

# Default command runs the full system test
CMD ["/usr/local/bin/run-migration-test.sh"]
