# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOENV=$(GOCMD) env

GOPATH=$(shell $(GOENV) GOPATH)
GOEXE=$(shell $(GOENV) GOEXE)

LDFLAGS=-s -w

# Binary names
APP_NAME=aiub-notice
LAUNCHER_NAME=$(APP_NAME)-launcher

# Source directory
APP_SRC=$(CURDIR)/cmd/$(APP_NAME)
LAUNCHER_SRC=$(CURDIR)/cmd/$(LAUNCHER_NAME)

# Build directory
BUILD_DIR=$(CURDIR)/bin

.PHONY: all build clean deps dev help install test uninstall

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /;s/: / - /'

## all: All-in-one target to clean, build, and update dependencies
all: clean deps build

## Clean: Remove build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

## dev: Build the application (development mode)
dev:
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)$(GOEXE) $(APP_SRC)
	$(GOBUILD) -ldflags="-H=windowsgui" -o $(BUILD_DIR)/$(LAUNCHER_NAME)$(GOEXE) $(LAUNCHER_SRC)

## build: Build the application (production mode)
build:
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)$(GOEXE) $(APP_SRC)
	$(GOBUILD) -ldflags="$(LDFLAGS) -H=windowsgui" -o $(BUILD_DIR)/$(LAUNCHER_NAME)$(GOEXE) $(LAUNCHER_SRC)

## test: Run tests
test:
	$(GOTEST) ./...

## deps: Update dependencies
deps:
	$(GOMOD) tidy
	$(GOGET) -u ./...

## install: Install the application
install:
	$(GOINSTALL) -ldflags="$(LDFLAGS)" $(APP_SRC)
	$(GOINSTALL) -ldflags="$(LDFLAGS) -H=windowsgui" $(LAUNCHER_SRC)

## uninstall: Uninstall the application
uninstall: _autostart_disable
	@read -p "Are you sure you want to uninstall $(APP_NAME)? (y/n): " confirm; \
		if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
			echo "Uninstalling $(APP_NAME)..."; \
		else \
			echo "Uninstallation cancelled."; \
			exit 1; \
		fi

	rm -f $(GOPATH)/bin/$(APP_NAME)$(GOEXE)
	rm -f $(GOPATH)/bin/$(LAUNCHER_NAME)$(GOEXE)

# autostart disable helper
_autostart_disable:
	@read -p "Remove autostart entry if exists? (y/n): " choice; \
		if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
			$(GOPATH)/bin/$(APP_NAME)$(GOEXE) autostart --disable || @echo "No autostart entry found or error occurred."; \
		else \
			echo "Skipping autostart entry removal."; \
		fi
