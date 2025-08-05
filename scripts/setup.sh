#!/bin/bash

echo "🚀 Setting up Network Sniffer Development Environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Python3 is not installed. Please install Python 3.9+ first."
    exit 1
fi

echo "✅ Go and Python are installed"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download

# Install pre-commit
echo "🔧 Installing pre-commit..."
pip install pre-commit

# Install git hooks
echo "🪝 Installing pre-commit hooks..."
pre-commit install

# Build application
echo "🔨 Building application..."
go build -o bin/network-sniffer cmd/server/main.go

# Generate swagger docs
echo "📚 Generating swagger documentation..."
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs

echo "✅ Setup complete!"
echo ""
echo "🎯 Next steps:"
echo "  1. Run the application: ./bin/network-sniffer"
echo "  2. Test pre-commit: pre-commit run --all-files"
echo "  3. View API docs: http://localhost:8080/swagger/index.html"
echo ""
echo "🪝 Pre-commit hooks are now active and will run on every commit!"
