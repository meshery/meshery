# Installation and Local Development

This guide helps you set up dependencies and run Meshery locally (server + UI).

## Prerequisites

- Git
- Docker (optional, for containerized workflows)
- Node.js 18 or 20 with npm
- Go 1.24.x

## Setup Steps

1. Install Node.js 18 or 20
   - Linux (NodeSource):
     ```bash
     curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
     sudo apt-get install -y nodejs
     ```
   - Or use `nvm` and select Node 18/20.

2. Install Go 1.24.x
   - Using tarball (no sudo):
     ```bash
     curl -fsSLo /tmp/go1.24.5.linux-amd64.tar.gz https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
     rm -rf "$HOME/.local/go"
     mkdir -p "$HOME/.local"
     tar -C "$HOME/.local" -xzf /tmp/go1.24.5.linux-amd64.tar.gz
     echo 'export PATH=$HOME/.local/go/bin:$PATH' >> ~/.bashrc
     export PATH=$HOME/.local/go/bin:$PATH
     go version
     ```

3. Install UI dependencies
   ```bash
   cd ui && npm ci
   cd ../provider-ui && npm ci
   ```

## Running Locally

Two processes are needed: the Go server and the Next.js UI. The UI proxies API calls to the server on port 9081.

1. Start the server
   ```bash
   export PATH=$HOME/.local/go/bin:$PATH
   make server PORT=9081
   ```

2. Start the UI (in another terminal)
   ```bash
   cd ui
   NODE_ENV=development PORT=3000 npm run dev
   ```

3. Open the app
   - UI: http://localhost:3000
   - Server (API): http://localhost:9081

## Notes

- The repository `Makefile` has additional targets, e.g., `ui`, `ui-build`, `server-without-operator`, etc.
- Ensure Node is v18 or v20 to avoid engine warnings.
- If using Docker, see `install/README.md` and root `README.md` for platform-specific options.


