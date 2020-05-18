OUTDIR := _output

ARCH := amd64

BINNAME_DARWIN ?= mesheryctl-darwin-$(ARCH)
BINNAME_LINUX ?= mesheryctl-linux-$(ARCH)
BINNAME_WINDOWS ?= mesheryctl-windows-$(ARCH).exe

# go options

# Need to remove the DWARF debug information and the symbol information
LDFLAGS := -w -s

# Load the parameters to go compiler and print optimization decisions
GCFLAGS := 


# CLI

.PHONY: $(BINNAME_DARWIN)
$(BINNAME_DARWIN):
	@echo
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=darwin GO111MODULE=on go build -ldflags='$(LDFLAGS)' -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_DARWIN) ./cmd/mesheryctl/main.go

.PHONY: $(BINNAME_LINUX)
$(BINNAME_LINUX):
	@echo
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=linux GO111MODULE=on go build -ldflags='$(LDFLAGS)' -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_LINUX) ./cmd/mesheryctl/main.go

.PHONY: $(BINNAME_WINDOWS)
$(BINNAME_WINDOWS):
	@echo
	CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=windows GO111MODULE=on go build -ldflags='$(LDFLAGS)' -gcflags='$(GCFLAGS)' -o $(OUTDIR)/$(BINNAME_WINDOWS) ./main.go
