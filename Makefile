# Set the name of your project
PROJECT_NAME := galera-ips

# Set the platforms to build for
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

# Set the flags to use for building
BUILD_FLAGS := -ldflags="-s -w" -tags=netgo -installsuffix netgo -trimpath

# Set the output directory
OUTPUT_DIR := bin

default: all

.PHONY: all
## Build the binary for all platforms
all: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) go build $(BUILD_FLAGS) -o $(OUTPUT_DIR)/$(PROJECT_NAME)_$(word 1,$(subst /, ,$@))_$(word 2,$(subst /, ,$@))$(if $(filter windows,$(word 1,$(subst /, ,$@))),.exe,) .

.PHONY: clean
## Remove the previous build
clean:
	rm -rf $(OUTPUT_DIR)/*

.PHONY: help
## This help screen
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-_0-9%:\\]+/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = $$1; \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			gsub("\\\\", "", helpCommand); \
			gsub(":+$$", "", helpCommand); \
			printf "  \x1b[32;01m%-20s\x1b[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"
