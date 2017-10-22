export VERSION=`git describe --abbrev=6 --dirty --always --tags`
export GOBIN=$(PWD)/bin
export CGO_ENABLED?=1
export installsuffix=cgo
#export GOFLAGS=-w -s

## keel ##
.PHONY: build
build: $(shell find . -name '*.go') .git
	go install -ldflags "-X main.version=${VERSION}" github.com/todaychiji/ha/cli/qi

.PHONY: run
run: build
	./bin/qi --version


.PHONY: server
server: build
	./bin/qi gateway

.PHONY: test
test:
	go test -v github.com/todaychiji/ha/...
