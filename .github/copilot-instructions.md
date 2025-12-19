# Copilot Instructions for Meshery

## Project Overview

Meshery is a cloud native manager that enables the design and management of Kubernetes-based infrastructure and applications. It is a CNCF project that supports 300+ integrations and provides visual and collaborative GitOps capabilities.

## Repository Structure

- **`/server`** - Meshery Server (Go/Golang backend)
- **`/ui`** - Meshery UI (React/Next.js frontend)
- **`/mesheryctl`** - Meshery CLI (Go/Golang command-line tool using Cobra framework)
- **`/docs`** - Documentation (Jekyll-based static site)
- **`/install`** - Installation scripts and Docker configurations

## Tech Stack

### Backend (Server & CLI)
- **Language**: Go 1.25.5
- **Framework**: Cobra (CLI), Go modules
- **Key Libraries**: MeshKit, Viper (config), gqlgen (GraphQL)

### Frontend (UI)
- **Framework**: Next.js with React
- **Design System**: Material UI (MUI)
- **State Management**: Redux Toolkit
- **API Client**: Relay (GraphQL), REST
- **Node Version**: 20 LTS

## Code Style and Conventions

### Go Code
- Follow standard Go conventions and formatting (`gofmt`)
- Use golangci-lint for linting (see `.golangci.yml`)
- Error handling should use MeshKit's error utilities
- Tests use Go's standard testing library

### JavaScript/React Code
- ESLint with Prettier for formatting (see `ui/.eslintrc.js`)
- Use functional components with hooks
- Follow Material UI patterns for styling
- Schema-driven UI development using JSON schemas

### Commit Messages
- Use descriptive PR titles with `[<component-name>]` prepended
- Sign commits with DCO (Developer Certificate of Origin) using `-s` flag
- Reference issue numbers in commits when applicable

## Building and Testing

### Server
```bash
make server          # Run Meshery Server locally (port 9081)
make build-server    # Build server binary
```

### UI
```bash
make ui-setup        # Install UI dependencies
make ui-build        # Build and export UI
make ui              # Run UI development server (port 3000)
```

### CLI (mesheryctl)
```bash
cd mesheryctl && make    # Build mesheryctl binary
go test --short ./...    # Run unit tests
go test -run Integration ./...  # Run integration tests
```

### Docker
```bash
make docker-build    # Build Docker container
```

## Common Development Tasks

### Adding New mesheryctl Commands
1. Follow the [Meshery CLI Style Guide](https://docs.meshery.io/project/contributing/contributing-cli-guide)
2. Add commands under `/mesheryctl/internal/cli/root/`
3. Use Cobra framework conventions
4. Include documentation in the command's Go file (Long, Short, Example fields)

### UI Component Development
1. Read the [Schema-Driven UI Development Guide](https://docs.meshery.io/project/contributing/contributing-ui-schemas)
2. Components go in `/ui/components/`
3. Use Sistent design system components from `@sistent/sistent` (Meshery's design system built on top of Material UI)
4. For components not available in Sistent, fall back to Material UI (MUI) components

### API Changes
- REST API: Available at `<hostname>:9081/api/`
- GraphQL API: Available at `<hostname>:9081/api/graphql/query`
- See [API Overview](https://docs.meshery.io/extensibility/api)

## Documentation

- Contributing guides: https://docs.meshery.io/project/contributing
- CLI documentation is auto-generated from Cobra command definitions
- Do not manually modify generated docs in `/docs` - edit source files in `/mesheryctl`

## Important Files

- `Makefile` - Main build targets
- `install/Makefile.core.mk` - Core build definitions
- `.golangci.yml` - Go linter configuration
- `ui/.eslintrc.js` - JavaScript linter configuration
- `ui/package.json` - UI dependencies and scripts

## Resources

- [Meshery Documentation](https://docs.meshery.io)
- [Contributing Guide](https://docs.meshery.io/project/contributing)
- [Community Slack](https://slack.meshery.io)
- [Meshery Architecture](https://docs.meshery.io/concepts/architecture)
