# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOENV=$(GOCMD) env

# Environment variables
GOPATH=$(shell $(GOENV) GOPATH)
GOEXE=$(shell $(GOENV) GOEXE)

# Project parameters
PKG=github.com/AtifChy/aiub-notice
APP_NAME=aiub-notice
LAUNCHER_NAME=$(APP_NAME)-launcher
VERSION=$(shell git describe --tags --always --dirty --long)

# Linker flags
LDFLAGS=-s -w -X $(PKG)/internal/common.Version=$(VERSION)

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
	$(GOBUILD) -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)$(GOEXE) $(APP_SRC)
	$(GOBUILD) -trimpath -ldflags="$(LDFLAGS) -H=windowsgui" -o $(BUILD_DIR)/$(LAUNCHER_NAME)$(GOEXE) $(LAUNCHER_SRC)

## test: Run tests
test:
	$(GOTEST) ./...

## deps: Update dependencies
deps:
	$(GOMOD) tidy
	$(GOGET) -u ./...

## install: Install the application
install:
	$(GOINSTALL) -trimpath -ldflags="$(LDFLAGS)" $(APP_SRC)
	$(GOINSTALL) -trimpath -ldflags="$(LDFLAGS) -H=windowsgui" $(LAUNCHER_SRC)
	@echo ""
	@echo "Installation complete!"
	@echo "Binaries installed to: $(GOPATH)/bin"
	@echo "Add $(GOPATH)/bin to your system PATH if not already done."

## uninstall: Uninstall the application
uninstall: _autostart_disable _aumid_deregister
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

# aumid deregister helper
_aumid_deregister:
	@read -p "Deregister AUMID if exists? (y/n): " choice; \
		if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
			$(GOPATH)/bin/$(APP_NAME)$(GOEXE) aumid --deregister || @echo "No AUMID found or error occurred."; \
		else \
			echo "Skipping AUMID deregistration."; \
		fi
