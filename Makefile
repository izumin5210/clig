PATH := ${PWD}/bin:${PATH}
export PATH

.DEFAULT_GOAL := all

REVISION ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date +'%Y-%m-%dT%H:%M:%SZ')

GO_BUILD_FLAGS := -v
GO_TEST_FLAGS := -v -timeout 2m
GO_COVER_FLAGS := -coverprofile coverage.txt -covermode atomic
SRC_FILES := $(shell go list -f '{{range .GoFiles}}{{printf "%s/%s\n" $$.Dir .}}{{end}}' ./...)

XC_ARCH := 386 amd64
XC_OS := darwin linux windows


#  App
#----------------------------------------------------------------
BIN_DIR := ./bin
OUT_DIR := ./dist
GENERATED_BINS :=
PACKAGES :=

define cmd-tmpl

$(eval NAME := $(notdir $(1)))
$(eval OUT := $(addprefix $(BIN_DIR)/,$(NAME)))
$(eval LDFLAGS := -ldflags "-X main.revision=$(REVISION) -X main.buildDate=$(BUILD_DATE)")

$(OUT): $(SRC_FILES)
	go build $(GO_BUILD_FLAGS) $(LDFLAGS) -o $(OUT) $(1)

.PHONY: $(NAME)
$(NAME): $(OUT)

.PHONY: $(NAME)-package
$(NAME)-package: $(NAME)
	gox \
		$(LDFLAGS) \
		-os="$(XC_OS)" \
		-arch="$(XC_ARCH)" \
		-output="$(OUT_DIR)/$(NAME)_{{.OS}}_{{.Arch}}" \
		$(1)

$(eval GENERATED_BINS += $(OUT))
$(eval PACKAGES += $(NAME)-package)

endef

$(foreach src,$(wildcard ./cmd/*),$(eval $(call cmd-tmpl,$(src))))


#  Commands
#----------------------------------------------------------------
.PHONY: tools
tools:
	go genereate ./tools.go

.PHONY: all
all: $(GENERATED_BINS)

.PHONY: packages
packages: $(PACKAGES)

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/*

.PHONY: gen
gen: tools
	go generate ./...

.PHONY: test
test:
	go test $(GO_TEST_FLAGS) ./...

.PHONY: cover
cover:
	go test $(GO_TEST_FLAGS) $(GO_COVER_FLAGS) ./...

.PHONY: lint
lint: ./bin/reviewdog ./bin/golangci-lint
ifdef CI
	reviewdog -reporter=github-pr-review
else
	reviewdog -diff="git diff master"
endif

# linters
bin/reviewdog:
	curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh | sh -s -- -b ./bin v0.9.12

bin/golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ./bin v1.17.1
