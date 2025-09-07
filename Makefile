# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v
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

# path separator
ifeq ($(OS), Windows_NT)
	SEP = \\
else
	SEP = /
endif

.PHONY: all build clean deps dev help install install-all test uninstall uninstall-all

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
	$(GOBUILD) -o $(BUILD_DIR)/debug/$(APP_NAME)$(GOEXE) $(APP_SRC)
	$(GOBUILD) -ldflags="-H=windowsgui" -o $(BUILD_DIR)/debug/$(LAUNCHER_NAME)$(GOEXE) $(LAUNCHER_SRC)

## build: Build the application (production mode)
build:
	$(GOBUILD) -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/release/$(APP_NAME)$(GOEXE) $(APP_SRC)
	$(GOBUILD) -trimpath -ldflags="$(LDFLAGS) -H=windowsgui" -o $(BUILD_DIR)/release/$(LAUNCHER_NAME)$(GOEXE) $(LAUNCHER_SRC)

## test: Run tests
test:
	$(GOTEST) ./...

## deps: Update dependencies
deps:
	$(GOGET) -u ./...
	$(GOMOD) tidy

## install-all: Install the application and optionally enable autostart and register AUMID
install-all: install _autostart_enable _aumid_register
	@echo ""
	@echo "All setup complete!"
	@echo ""
	@echo "Installation complete!"
	@echo "Binaries installed to: $(GOPATH)$(SEP)bin"
	@echo ""
	@echo "Add $(GOPATH)$(SEP)bin to your system PATH if not already done."
	@echo "Run the following command in PowerShell to add it to your user path:"
	@echo "[Environment]::SetEnvironmentVariable(\"PATH\", \"\$$env:PATH;$(GOPATH)$(SEP)bin\", [EnvironmentVariableTarget]::User)"

## install: Install the application
install:
	$(GOINSTALL) -trimpath -ldflags="$(LDFLAGS)" $(APP_SRC)
	$(GOINSTALL) -trimpath -ldflags="$(LDFLAGS) -H=windowsgui" $(LAUNCHER_SRC)

# autostart enable helper
_autostart_enable:
	@read -p "Do you want to enable autostart? [recommended] (y/n): " choice; \
		if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
			$(GOPATH)/bin/$(APP_NAME)$(GOEXE) autostart --enable || @echo "Failed to enable autostart."; \
		else \
			echo "Skipping autostart enable."; \
		fi

# aumid register helper
_aumid_register:
	@read -p "Do you want to register AUMID? (y/n): " choice; \
		if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
			$(GOPATH)/bin/$(APP_NAME)$(GOEXE) aumid --register || @echo "Failed to register AUMID."; \
		else \
			echo "Skipping AUMID registration."; \
		fi

## uninstall-all: Uninstall the application and optionally disable autostart and unregister AUMID
uninstall-all: uninstall _autostart_disable _aumid_unregister
	@echo ""
	@echo "Uninstallation complete!"

## uninstall: Uninstall the application
uninstall:
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

# aumid unregister helper
_aumid_unregister:
	@read -p "Unregister AUMID if exists? (y/n): " choice; \
		if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
			$(GOPATH)/bin/$(APP_NAME)$(GOEXE) aumid --unregister || @echo "No AUMID found or error occurred."; \
		else \
			echo "Skipping AUMID unregistration."; \
		fi
