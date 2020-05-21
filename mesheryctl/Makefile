make:
	go build -o mesheryctl cmd/mesheryctl/main.go
forwin:
	go build -o mesheryctl.exe cmd/mesheryctl/main.go

# ---------------------------------------------------------------------------------
# build for mutiple platform

TAGS :=
TARGETS :=darwin/amd64 linux/amd64 windows/amd64
BINNAME ?=mesheryctl

GOPATH = $(shell go env GOPATH)
GOX = $(GOPATH)/bin/gox

# go options

# Need to remove the DWARF debug information and the symbol information
LDFLAGS := -w -s
# Load the parameters to go compiler and print optimization decisions
GCFLAGS := 

$(GOX):
	(cd /; GO111MODULE=on go get -u github.com/mitchellh/gox)


.PHONY: build-multiple
build-multiple: LDFLAGS += -extldflags "-static"
build-multiple: $(GOX)
	GO111MODULE=on CGO_ENABLED=0 $(GOX) -parallel=3 -output="_dist/{{.OS}}-{{.Arch}}/$(BINNAME)" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/mesheryctl/

.PHONY: clean
clean:
	@rm -rf ./_dist