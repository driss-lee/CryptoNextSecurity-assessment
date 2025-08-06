# Network Sniffer Service

A Go-based network packet sniffing service that simulates packet capture and provides a REST API for querying packet data.

## 🚀 Features

- **Packet Simulation**: Generates realistic network packets with various protocols
- **REST API**: HTTP endpoints for querying packet data with filtering
- **Swagger Documentation**: Auto-generated API documentation
- **Docker Support**: Containerized deployment
- **Environment Configuration**: Support for development and production environments
- **Code Quality**: Pre-commit hooks for consistent code formatting and quality
- **CI/CD Pipeline**: Automated testing, building, and deployment
- **Live Deployment**: Available at https://cryptonextsecurity-assessment.onrender.com/

## 📋 Prerequisites

- Go 1.21+
- Docker (optional)
- Git
- Python 3.9+ (for pre-commit hooks)

## 🛠️ Installation

### Quick Setup (Recommended)

For a quick setup with all dependencies and pre-commit hooks:

```bash
# Clone the repository
git clone <repository-url>
cd network-sniffer

# Run the setup script
chmod +x scripts/setup.sh
./scripts/setup.sh
```

The setup script will:
- ✅ Check for required dependencies (Go 1.21+, Python 3.9+)
- 📦 Install Go dependencies
- 🔧 Install pre-commit hooks
- 🔨 Build the application
- 📚 Generate swagger documentation

### Manual Installation

If you prefer manual setup or the setup script fails:

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd network-sniffer
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Setup pre-commit hooks (Recommended):**
   ```bash
   # Install pre-commit
   pip install pre-commit

   # Install git hooks
   pre-commit install

   # Run against all files (optional)
   pre-commit run --all-files
   ```

4. **Build the application:**
   ```bash
   go build -o bin/network-sniffer cmd/server/main.go
   ```

5. **Run the service:**
   ```bash
   ./bin/network-sniffer
   ```

### Using Task Runner

```bash
# Install Task
go install github.com/go-task/task/v3/cmd/task@latest

# Available tasks
task build      # Build the application
task run        # Run the application
task test       # Run tests
task clean      # Clean build artifacts
task docker     # Build Docker image
task swagger    # Generate swagger docs
```

## 🐳 Docker

### Build and Run

```bash
# Build image
docker build -t network-sniffer .

# Run with environment file (Recommended)
docker run -p 8080:8080 --env-file .env.development network-sniffer

# Run with environment variables
docker run -p 8080:8080 \
  -e STORAGE_MAX_SIZE=5000 \
  -e SNIFFING_INTERVAL=2s \
  -e SERVER_PORT=8080 \
  network-sniffer
```

### Environment-Specific Docker Runs

```bash
# Development (uses .env.development)
docker run -p 8080:8080 --env-file .env.development network-sniffer

# Production simulation (uses .env.production)
ENV=production docker run -p 8080:8080 --env-file .env.production network-sniffer

# High-performance settings
docker run -p 8080:8080 \
  -e STORAGE_MAX_SIZE=10000 \
  -e SNIFFING_INTERVAL=1s \
  -e SERVER_SHUTDOWN_TIMEOUT=60s \
  network-sniffer
```

## 🚀 Deployment

### Render (Free Tier)

The project includes a `render.yaml` file for automatic deployment configuration.

1. **Create Render Account:**
   - Sign up at [render.com](https://render.com)
   - Connect your GitHub repository

2. **Automatic Deployment:**
   - Render will automatically deploy on every push to main branch
   - Free tier includes 750 hours/month
   - Environment variables are pre-configured in `render.yaml`

**Environment Configuration:**
- Environment variables are loaded from the `.env.development` file (development) or `.env.production` file (production)
- `PORT=8080` - Render's port assignment (set automatically by Render)
- All other configuration comes from the appropriate environment file

**Recommended Production Setup:**
Add these to your Render dashboard under Environment tab:
```
ENV=production
STORAGE_MAX_SIZE=5000
SNIFFING_INTERVAL=2s
SERVER_PORT=8080
SERVER_SHUTDOWN_TIMEOUT=60s
```

### GitHub Actions CI/CD

The project includes a comprehensive CI/CD pipeline with the following jobs:

#### Code Quality Job
- **Runs pre-commit checks** on every push and pull request
- **Validates code quality** and formatting
- **Ensures consistency** across all contributions
- **Checks include**: file formatting, YAML/JSON validation, Go formatting, imports, tests, build validation

#### Build Job
- **Depends on code-quality** job completion
- **Builds the application** and Docker image
- **Generates swagger documentation**
- **Uploads build artifacts** for testing

#### Test Job
- **Depends on build** job completion
- **Runs unit tests** and integration tests
- **Tests Docker image** with custom environment variables
- **Downloads build artifacts** for testing

#### Deploy Job
- **Depends on test** job completion
- **Deploys to Render** only when merged to main branch
- **Runs on pull requests** and main branch pushes

#### Workflow Triggers

| Event | Code Quality | Build | Test | Deploy |
|-------|-------------|-------|------|--------|
| **Push to develop** | ✅ | ✅ | ✅ | ❌ |
| **Pull Request to main** | ✅ | ✅ | ✅ | ✅ |
| **Merge to main** | ✅ | ✅ | ✅ | ✅ |

> **Note**: Direct pushes to main are restricted. All changes must go through pull requests.

#### Setup GitHub Secrets

Add these secrets to your GitHub repository:

- `RENDER_TOKEN` - Your Render API token
- `RENDER_SERVICE_ID` - Your Render service ID

## 🔧 Code Quality

### Pre-commit Hooks

The project uses pre-commit hooks to ensure code quality:

#### File Hygiene
- **trailing-whitespace**: Removes trailing whitespace
- **end-of-file-fixer**: Ensures files end with newline
- **check-yaml**: Validates YAML files
- **check-json**: Validates JSON files
- **check-added-large-files**: Prevents large files from being committed
- **check-merge-conflict**: Detects merge conflict markers
- **check-case-conflict**: Detects case conflicts in filenames

#### Go Code Quality
- **go-fmt**: Formats Go code
- **go-imports**: Organizes Go imports
- **go-unit-tests**: Runs Go tests
- **go-build**: Validates Go build
- **go-mod-tidy**: Manages Go dependencies

#### Manual Testing
```bash
# Run all pre-commit hooks
pre-commit run --all-files

# Run specific hook
pre-commit run go-fmt

# Run on staged files only
pre-commit run
```

#### Git Integration
```bash
# Install git hooks (automatic on commit)
pre-commit install

# Uninstall git hooks
pre-commit uninstall
```

## 🏗️ Architecture

```
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── config/         # Configuration management
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   └── storage/        # Data storage layer
├── pkg/
│   └── sniffing/      # Packet sniffing simulation
├── docs/              # Generated swagger documentation
├── bin/               # Build artifacts (gitignored)
├── .github/           # GitHub Actions workflows
├── .pre-commit-config.yaml  # Pre-commit configuration
├── .env.development   # Development environment
└── .env.production    # Production environment
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/storage
```

### API Testing

#### Local Testing
```bash
# Test packets endpoint
curl http://localhost:8080/api/v1/packets

# Test with filters
curl "http://localhost:8080/api/v1/packets?protocol=TCP&limit=5"

# Test swagger docs
curl http://localhost:8080/swagger/doc.json
```

#### Live Deployment Testing
```bash
# Test packets endpoint
curl https://cryptonextsecurity-assessment.onrender.com/api/v1/packets

# Test with filters
curl "https://cryptonextsecurity-assessment.onrender.com/api/v1/packets?protocol=TCP&limit=5"

# Access Swagger documentation
# https://cryptonextsecurity-assessment.onrender.com/swagger/index.html
```

### Logs

The application logs to stdout with basic information:
- Server startup/shutdown
- Packet sniffing start/stop
- Error messages

## 🔧 Configuration

### Environment Variables

The application supports environment variables for deployment customization:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `ENV` | Environment mode | `development` | `production` |
| `STORAGE_MAX_SIZE` | Maximum packets in memory | `1000` | `5000` |
| `SNIFFING_INTERVAL` | Packet generation interval | `5s` | `2s` |
| `SERVER_PORT` | HTTP server port | `8080` | `3000` |
| `SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `30s` | `60s` |

### Environment Files

For deployment flexibility, environment files are available:

#### Development Configuration (`.env.development`)
```bash
STORAGE_MAX_SIZE=100
SNIFFING_INTERVAL=10s
SERVER_PORT=8080
SERVER_SHUTDOWN_TIMEOUT=30s
```

#### Production Configuration (`.env.production`)
```bash
STORAGE_MAX_SIZE=5000
SNIFFING_INTERVAL=2s
SERVER_PORT=8080
SERVER_SHUTDOWN_TIMEOUT=60s
```

### Local Development

#### Using Environment Files
```bash
# Development (uses .env.development)
./network-sniffer

# Production simulation (uses .env.production)
ENV=production ./network-sniffer
```

#### Using Environment Variables
```bash
# Override specific values
STORAGE_MAX_SIZE=5000 SNIFFING_INTERVAL=2s ./network-sniffer

# Full custom configuration
STORAGE_MAX_SIZE=10000 \
SNIFFING_INTERVAL=1s \
SERVER_PORT=3000 \
SERVER_SHUTDOWN_TIMEOUT=60s \
./network-sniffer
```

### Docker Development

#### Using Environment Files
```bash
# Development
docker run -p 8080:8080 --env-file .env.development network-sniffer

# Production simulation
ENV=production docker run -p 8080:8080 --env-file .env.production network-sniffer
```

#### Using Environment Variables
```bash
# Custom settings
docker run -p 8080:8080 \
  -e STORAGE_MAX_SIZE=10000 \
  -e SNIFFING_INTERVAL=1s \
  -e SERVER_PORT=3000 \
  network-sniffer
```

### CI/CD (GitHub Actions)

The CI/CD pipeline tests environment variable functionality:

```yaml
# Tests Docker image with custom environment variables
- name: Test Docker image with custom environment variables
  run: |
    timeout 10s docker run --rm \
      -e STORAGE_MAX_SIZE=100 \
      -e SNIFFING_INTERVAL=1s \
      -e SERVER_PORT=8080 \
      -e SERVER_SHUTDOWN_TIMEOUT=10s \
      network-sniffer || true
```

### Production Deployment (Render)

#### Option A: Render Dashboard (Recommended)
1. Go to your Render service dashboard
2. Navigate to **Environment** tab
3. Add environment variables:
   ```
   ENV=production
   STORAGE_MAX_SIZE=5000
   SNIFFING_INTERVAL=2s
   SERVER_PORT=8080
   SERVER_SHUTDOWN_TIMEOUT=60s
   ```

#### Option B: Git Files Only
1. Add only this to Render dashboard:
   ```
   ENV=production
   ```
2. The app will use `.env.production` from your git repository

#### Option C: Hybrid Approach
1. Keep `.env.production` in git for defaults
2. Override specific values in Render dashboard:
   ```
   ENV=production
   STORAGE_MAX_SIZE=10000  # Override default
   SNIFFING_INTERVAL=1s    # Override default
   ```

### Environment File Structure

```
├── .env.development    # Development configuration (committed to git)
├── .env.production     # Production configuration (committed to git)
├── .env.local          # Local overrides (gitignored)
└── .env.*.local        # Environment-specific local overrides (gitignored)
```

## 🤝 Contributing

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes** (pre-commit hooks will run automatically)
4. **Commit your changes:**
   ```bash
   git add .
   git commit -m 'Add amazing feature'
   ```
5. **Push to the branch:**
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Open a Pull Request**

### Development Guidelines

- **Follow Go best practices**
- **Add tests for new features**
- **Update swagger documentation**
- **Keep commits atomic and well-described**
- **Ensure pre-commit hooks pass**
- **All code must pass CI/CD pipeline**

### Branch Protection

The main branch is protected with the following rules:
- **Requires pull requests** for all changes
- **Requires CI/CD to pass** before merge
- **Requires code review** before merge
- **No direct pushes** to main branch

## 🆘 Support

- **API Documentation**:
  - Local: `http://localhost:8080/swagger/index.html`
  - Live: `https://cryptonextsecurity-assessment.onrender.com/swagger/index.html`
- **Live API**: `https://cryptonextsecurity-assessment.onrender.com/api/v1/packets`
- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
