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

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -installsuffix cgo -a -o bin/${SERVICE}d ${PACKAGE_CMD_DAEMON}
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -installsuffix cgo -a -o bin/${SERVICE}stress ${PACKAGE_CMD_STRESS}

install:
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_DAEMON}
	go install -ldflags "${LDFLAGS}" ${PACKAGE_CMD_STRESS}

test:
	scripts/test.sh
	go tool cover -func=coverage.txt | tail -n 1

cover:
	go tool cover -html=coverage.txt

get:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

publish: build
	docker build \
		--build-arg VCS_REF=${VCS_REF} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t piotrkowalczuk/${SERVICE}:${TAG} .
	docker push piotrkowalczuk/${SERVICE}:${TAG}