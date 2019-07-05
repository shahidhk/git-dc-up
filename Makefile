DOCKER_REGISTRY   ?= docker.io
IMAGE_PREFIX      ?= shahidh
SHORT_NAME        ?= git-dc-up

# build options
GO        ?= go
TAGS      :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

.PHONY: all
all: build

.PHONY: build
build:
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)'

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

include versioning.mk
