APP_NAME := deploykf

BINDIR       := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin

GIT_COMMIT     := $(shell git rev-parse HEAD)
GIT_TAG        := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_TREE_STATE := $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
ifeq ($(GIT_TAG),)
    GIT_TAG := v0.0.0
endif
BINARY_VERSION ?= ${GIT_TAG}

LDFLAGS     := -w -s
LDFLAGS     += -X github.com/deployKF/cli/internal/version.version=${BINARY_VERSION}
LDFLAGS     += -X github.com/deployKF/cli/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS     += -X github.com/deployKF/cli/internal/version.gitTreeState=${GIT_TREE_STATE}

TARGETS := \
	linux-amd64 \
	linux-arm64 \
	darwin-amd64 \
	darwin-arm64 \
	windows-amd64

# ------------------------------------------------------------------------------
#  checks

.PHONY: check-golangci-lint
check-golangci-lint:
ifeq (, $(shell which golangci-lint))
	$(error "golangci-lint is not installed. Please install it using the following command: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2")
endif

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build:
	@echo "********** building $(APP_NAME) for $(shell go env GOOS)/$(shell go env GOARCH) **********"
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(APP_NAME)'

# ------------------------------------------------------------------------------
#  build-all

.PHONY: build-all
build-all: $(addprefix build-, $(TARGETS))
build-%:
	$(eval os_arch := $(subst -, ,$*))
	$(eval GOOS := $(word 1, $(os_arch)))
	$(eval GOARCH := $(word 2, $(os_arch)))
	$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))

	@echo "********** building $(APP_NAME) for $(GOOS)/$(GOARCH) **********"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o '$(BINDIR)/$(APP_NAME)-$(GOOS)-$(GOARCH)$(EXT)'

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@echo "********** installing into $(INSTALL_PATH)/$(APP_NAME) **********"
	@install '$(BINDIR)/$(APP_NAME)' '$(INSTALL_PATH)/$(APP_NAME)'

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test:
	@echo "********** running tests **********"
	@go test -v ./...

# ------------------------------------------------------------------------------
#  lint

.PHONY: lint
lint: check-golangci-lint
	@echo "********** running golangci-lint **********"
	@golangci-lint run ./...

# ------------------------------------------------------------------------------
#  lint-fix

.PHONY: lint-fix
lint-fix: check-golangci-lint
	@echo "********** running gofmt and goimports **********"
	@find . -name '*.go' -type f -not -path './vendor/*' -exec gofmt -s -w {} \;
	@goimports -local github.com/deployKF/cli -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

	@echo "********** running golangci-lint --fix **********"
	@golangci-lint run --fix ./...

# ------------------------------------------------------------------------------
#  clean

.PHONY: clean
clean:
	@echo "********** cleaning up build artifacts **********"
	@rm -rf '$(BINDIR)'
