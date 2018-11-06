SERVICE=mnemosyne
VERSION=$(shell git describe --tags --always --dirty)
ifeq ($(version),)
	TAG=${VERSION}
else
	TAG=$(version)
endif

PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
PACKAGE_CMD_STRESS=$(PACKAGE)/cmd/$(SERVICE)stress

LDFLAGS = -X 'main.version=$(VERSION)'

.PHONY:	all gen build install test cover get

all: get install

version:
	echo ${VERSION} > VERSION.txt

gen:
	./scripts/generate.sh
	./.circleci/scripts/generate.sh golang

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -installsuffix cgo -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -installsuffix cgo -a -o bin/${SERVICE}stress ${PACKAGE_CMD_STRESS}

install:
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_STRESS}

test:
	./.circleci/scripts/test.sh
	go tool cover -func=cover.out | tail -n 1

cover: test
	go tool cover -html=cover.out

get:
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	go get -u google.golang.org/grpc
	go get -u github.com/axw/gocov/gocov
	go get -u gotest.tools/gotestsum
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

publish: build
ifneq ($(skiplogin),true)
	docker login
endif
	docker build \
		--build-arg VCS_REF=${VCS_REF} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t piotrkowalczuk/${SERVICE}:${TAG} .
	docker push piotrkowalczuk/${SERVICE}:${TAG}
