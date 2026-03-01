---
layout: page
title: Development Environment Troubleshooting
permalink: guides/troubleshooting/development-environment
type: guides
category: troubleshooting
language: en
abstract: Comprehensive troubleshooting guide for common development environment setup issues when contributing to Meshery.
---

# Development Environment Troubleshooting

This guide helps you resolve common issues encountered when setting up your development environment for Meshery contributions.

## Go Environment Issues

### GOPATH Configuration Problems

**Problem:** `$GOPATH not set` or `cannot find package` errors

**Solution:**
```bash
# Check if GOPATH is set
echo $GOPATH

# Set GOPATH (add to your shell profile)
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Verify Go installation
go version
go env GOPATH
```

### Go Version Conflicts

**Problem:** Multiple Go versions causing conflicts

**Solution:**
```bash
# Check current Go version
go version

# Install specific Go version (if needed)
# Download from https://golang.org/dl/

# Update PATH to use correct version
export PATH=/usr/local/go/bin:$PATH

# Verify correct version is active
which go
go version
```

### Module Resolution Issues

**Problem:** `go: cannot find module` or dependency resolution errors

**Solution:**
```bash
# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Verify module integrity
go mod verify

# Tidy up dependencies
go mod tidy
```

## Node.js/npm Issues

### Version Compatibility Problems

**Problem:** Node.js version incompatibility (version 19+ not supported)

**Solution:**
```bash
# Check current Node.js version
node --version

# Install Node Version Manager (nvm)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Install and use compatible Node.js version
nvm install 18
nvm use 18
nvm alias default 18

# Verify version
node --version
```

### Package Installation Failures

**Problem:** `npm install` fails with permission or network errors

**Solution:**
```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules and package-lock.json
rm -rf node_modules package-lock.json

# Reinstall dependencies
npm install

# For permission issues on macOS/Linux
sudo chown -R $(whoami) ~/.npm
```

### Build Errors

**Problem:** `npm run build` or `make ui-build` fails

**Solution:**
```bash
# Check for TypeScript errors
npm run type-check

# Clear build cache
rm -rf .next out dist

# Rebuild from scratch
npm run clean
npm run build

# Check for memory issues
export NODE_OPTIONS="--max-old-space-size=4096"
npm run build
```

## Docker/Container Issues

### Docker Daemon Connection Problems

**Problem:** `Cannot connect to the Docker daemon`

**Solution:**
```bash
# Start Docker daemon (Linux)
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# For Docker Desktop (Windows/macOS)
# Ensure Docker Desktop is running

# Test Docker connection
docker ps
```

### Port Conflicts

**Problem:** `Port already in use` errors

**Solution:**
```bash
# Find process using port (e.g., 9081)
lsof -i :9081
netstat -tulpn | grep 9081

# Kill process using port
kill -9 <PID>

# Or use different port
export PORT=9082
make server
```

### Volume Mounting Issues

**Problem:** Docker volume mount failures

**Solution:**
```bash
# Check Docker volume permissions
ls -la /path/to/volume

# Fix permissions (Linux/macOS)
sudo chown -R $(whoami):$(whoami) /path/to/volume

# For Windows WSL2
# Ensure path is in WSL filesystem, not Windows
```

## Git/GitHub Issues

### Fork and Clone Problems

**Problem:** Unable to fork or clone repository

**Solution:**
```bash
# Verify GitHub authentication
ssh -T git@github.com

# Set up SSH key if needed
ssh-keygen -t ed25519 -C "your_email@example.com"
cat ~/.ssh/id_ed25519.pub
# Add to GitHub SSH keys

# Clone with SSH instead of HTTPS
git clone git@github.com:yourusername/meshery.git
```

### Remote Configuration

**Problem:** Incorrect remote configuration

**Solution:**
```bash
# Check current remotes
git remote -v

# Add upstream remote
git remote add upstream https://github.com/meshery/meshery.git

# Verify remotes
git remote -v

# Should show:
# origin    git@github.com:yourusername/meshery.git (fetch)
# origin    git@github.com:yourusername/meshery.git (push)
# upstream  https://github.com/meshery/meshery.git (fetch)
# upstream  https://github.com/meshery/meshery.git (push)
```

### DCO Sign-off Issues

**Problem:** Missing Developer Certificate of Origin (DCO) sign-off

**Solution:**
```bash
# Configure git with your details
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Sign-off commits automatically
git config --global alias.cm "commit -s -m"

# Amend last commit to add sign-off
git commit --amend --signoff

# For multiple commits, use interactive rebase
git rebase -i HEAD~n --signoff
```

## Platform-specific Issues

### Windows WSL Setup

**Problem:** WSL-related development issues

**Solution:**
```bash
# Update WSL to version 2
wsl --set-default-version 2
wsl --set-version Ubuntu-22.04 2

# Install required packages in WSL
sudo apt update
sudo apt install -y build-essential git curl

# Set up proper line endings
git config --global core.autocrlf input
git config --global core.eol lf

# Use WSL filesystem for better performance
cd /home/username/projects
```

### macOS Permission Issues

**Problem:** Permission denied errors on macOS

**Solution:**
```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) /usr/local/Homebrew

# Install Xcode command line tools
xcode-select --install

# Fix npm permissions
sudo chown -R $(whoami) ~/.npm
sudo chown -R $(whoami) /usr/local/lib/node_modules
```

### Linux Dependency Problems

**Problem:** Missing system dependencies

**Solution:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y build-essential git curl wget \
  ca-certificates gnupg lsb-release

# CentOS/RHEL/Fedora
sudo yum groupinstall -y "Development Tools"
sudo yum install -y git curl wget

# Install Docker (Ubuntu)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

## Common Build Errors

### Make Command Failures

**Problem:** `make` commands fail with various errors

**Solution:**
```bash
# Ensure make is installed
which make
sudo apt install make  # Ubuntu/Debian
brew install make      # macOS

# Clean build artifacts
make clean

# Check Makefile for specific targets
make help

# Run with verbose output
make server VERBOSE=1
```

### Dependency Resolution

**Problem:** Dependency conflicts or missing dependencies

**Solution:**
```bash
# Update all dependencies
go get -u ./...
npm update

# Clear all caches
go clean -cache -modcache -i -r
npm cache clean --force

# Reinstall from scratch
rm -rf node_modules vendor
go mod download
npm install
```

### Test Execution Problems

**Problem:** Tests fail to run or execute

**Solution:**
```bash
# Run tests with verbose output
go test -v ./...
npm test -- --verbose

# Run specific test
go test -v ./path/to/package -run TestName
npm test -- --testNamePattern="test name"

# Check test coverage
go test -cover ./...
npm test -- --coverage
```

## Getting Help

If you're still experiencing issues after trying these solutions:

1. **Search existing issues:** Check [GitHub Issues](https://github.com/meshery/meshery/issues) for similar problems
2. **Join the community:** Get help in [Meshery Slack](https://slack.meshery.io)
3. **Create an issue:** If you've found a new problem, [create an issue](https://github.com/meshery/meshery/issues/new/choose)
4. **Ask in discussions:** Use [GitHub Discussions](https://github.com/meshery/meshery/discussions) for questions

## Contributing Back

Found a solution not covered here? Please contribute back by:
- Opening a pull request to add your solution to this guide
- Sharing your experience in the community Slack
- Updating relevant documentation

Remember: Every contributor was once a beginner. Your questions and solutions help make Meshery more accessible to everyone! 🎈
