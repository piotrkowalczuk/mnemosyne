VERSION=$(shell git describe --tags --always --dirty)
SERVICE=mnemosyne

PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
PACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /mnemosynerpc| grep -v /mnemosynetest)

.PHONY:	all gen build install test cover get

all: get install

gen:
	@go generate .
	@go generate ./${SERVICE}d
	@go generate ./${SERVICE}rpc
	@ls -al ./${SERVICE}rpc | grep "pb.go"

build:
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}

install:
	@go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}

test:
	@touch .tmp/coverage.out
	@echo "mode: atomic" > .tmp/coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=.tmp/coverage.out -covermode=atomic $(pkg) || exit;\
		tail -n +2 .tmp/coverage.out >> .tmp/coverage-all.out \
	;)
	@go tool cover -func=.tmp/coverage-all.out | tail -n 1

cover: test
	@go tool cover -html=.tmp/coverage-all.out

get:
	@go get github.com/Masterminds/glide
	@glide --no-color install

publish:
	@docker build -t piotrkowalczuk/${SERVICE}:${VERSION} .
	@docker push piotrkowalczuk/${SERVICE}:${VERSION}
