export VERSION=`git describe --abbrev=6 --dirty --always --tags`
export GOBIN=$(PWD)/build
export CGO_ENABLED?=0
export installsuffix=cgo
#export GOFLAGS=-w -s

## keel ##
.PHONY: build
build: $(shell find . -name '*.go')
	rm build/up
	go install -ldflags "-X main.version=${VERSION}" github.com/todaychiji/ha/cli/up

.PHONY: run
run: build
	./build/up --version

.PHONY: test
test:
	go test -v github.com/DaoCloud/blackpearl/keel/cmd/keel/...
