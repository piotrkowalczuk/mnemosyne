VERSION?=$(shell git describe --tags --always --dirty)
SERVICE=mnemosyne

PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
PACKAGE_CMD_STRESS=$(PACKAGE)/cmd/$(SERVICE)stress

LDFLAGS = -X 'main.version=$(VERSION)'

.PHONY:	all gen build install test cover get

all: get install

version:
	@echo ${VERSION} > VERSION.txt

gen:
	@go generate .
	@go generate ./${SERVICE}d
	@go generate ./${SERVICE}rpc
	@ls -al ./${SERVICE}rpc | grep "pb.go"

build:
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -installsuffix cgo -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}
	@CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -installsuffix cgo -a -o bin/${SERVICE}stress ${PACKAGE_CMD_STRESS}

install:
	@go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}
	@go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_STRESS}

test:
	@scripts/test.sh
	@go tool cover -func=coverage.txt | tail -n 1

cover: test
	@go tool cover -html=coverage.txt

get:
	@go get github.com/Masterminds/glide
	@glide --no-color install

publish:
	@docker build \
		--build-arg VCS_REF=${VCS_REF} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t piotrkowalczuk/${SERVICE}:${VERSION} .
	@docker push piotrkowalczuk/${SERVICE}:${VERSION}