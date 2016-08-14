VERSION?=$(shell git describe --tags --always --dirty)
SERVICE=mnemosyne

PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d

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
	@./test.sh
	@go tool cover -func=coverage.txt | tail -n 1

cover: test
	@go tool cover -html=coverage.txt

get:
	@go get github.com/Masterminds/glide
	@glide --no-color install

publish:
	@docker build -t piotrkowalczuk/${SERVICE}:${VERSION} .
	@docker push piotrkowalczuk/${SERVICE}:${VERSION}
