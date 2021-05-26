PREFIX := github.com/kyoukaya/genshindaily
TARGETOS := linux
TARGETARCH := amd64
CGO_ENABLED := 0
GOFLAGS := -v
VERSION := git describe --tags
COMMIT := git rev-parse --short=8 HEAD
INT_PKG = $(PREFIX)/internal/genshindaily
LDFLAGS := -s -w -X '$(INT_PKG).Version=`$(VERSION)`' -X '$(INT_PKG).Date=`date`' -X '$(INT_PKG).Commit=`$(COMMIT)`'
RELEASE_DIR := dist
BIN := main

default: build

clean:
	rm -rf $(RELEASE_DIR) daily

release: build
	cd $(RELEASE_DIR) && zip -9 $(BIN).zip $(BIN)

build:
	go version
	GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) CGO_ENABLED=$(CGO_ENABLED) go build \
	$(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(BIN)

build-cmd:
	go version
	GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) CGO_ENABLED=$(CGO_ENABLED) go build \
	$(GOFLAGS) -ldflags "$(LDFLAGS)" -o ./daily ./cmd/dailycmd

test:
	mkdir -p .tmp
	go test -v -count 1 -coverprofile=.tmp/coverage.out ./... 
	go tool cover -func=.tmp/coverage.out

lint:
	golangci-lint version
	golangci-lint run -v ./...

format:
	goimports -local $(PREFIX) -w .
	gofmt -s -w .
