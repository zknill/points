PKGS := $(shell go list ./...)
FILES := $(shell find . -not -path "*vendor/*" -not -path "*.history/*" -name "*.go" | xargs -I % dirname % | sed 's/^.\///;s/[^.].*$$/&\/*.go/;s/^\.$$/*.go/' | sort -u)
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
GOIMPORTS := $(BIN_DIR)/goimports

default: test

.PHONY: test
test: lint
	go test $(PKGS) -timeout 1m

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter --vendor ./...  --deadline=180s

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: fmt
fmt: $(GOIMPORTS)
	goimports -w $(FILES)

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports