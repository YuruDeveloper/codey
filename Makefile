.PHONY: build run test clean lint fmt vet install help

# 변수 설정
BINARY_NAME=codey
BUILD_DIR=./bin
GO_FILES=$(shell find . -name '*.go' -type f)
MAIN_PATH=./cmd/main.go

# 기본 타겟
all: build

# 빌드
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 실행
run:
	@go run $(MAIN_PATH)

# 테스트
test:
	@echo "Running tests..."
	@go test -v ./...

# 테스트 커버리지
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 코드 포맷팅
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 코드 검사
vet:
	@echo "Running go vet..."
	@go vet ./...

# Lint (golangci-lint 필요)
lint:
	@echo "Running linter..."
	@golangci-lint run

# 의존성 다운로드
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# 설치
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(MAIN_PATH)

# 도움말
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run linter (requires golangci-lint)"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make install       - Install the binary"
	@echo "  make help          - Show this help message"