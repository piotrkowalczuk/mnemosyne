PROTOC=/usr/local/bin/protoc
SERVICE=mnemosyne
PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
BINARY=${SERVICE}d/${SERVICE}d

FLAGS=-host=$(MNEMOSYNE_HOST) \
          	    -port=$(MNEMOSYNE_PORT) \
          	    -subsystem=$(MNEMOSYNE_SUBSYSTEM) \
          	    -namespace=$(MNEMOSYNE_NAMESPACE) \
          	    -l.format=$(MNEMOSYNE_LOGGER_FORMAT) \
          	    -l.adapter=$(MNEMOSYNE_LOGGER_ADAPTER) \
          	    -l.level=$(MNEMOSYNE_LOGGER_LEVEL) \
          	    -m.engine=$(MNEMOSYNE_MONITORING_ENGINE) \
          	    -s.engine=$(MNEMOSYNE_STORAGE_ENGINE) \
          	    -sp.connectionstring=$(MNEMOSYNE_STORAGE_POSTGRES_CONNECTION_STRING) \
          	    -sp.tablename=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE_NAME) \
          	    -sp.retry=$(MNEMOSYNE_STORAGE_POSTGRES_RETRY)

.PHONY:	all proto build build-daemon run test test-unit test-postgres

all: proto build test run

proto:
	@${PROTOC} --proto_path=${GOPATH}/src --proto_path=. --go_out=plugins=grpc:. timestamp.proto ${SERVICE}.proto
	@ls -al | grep "pb.go"

build: build-daemon

build-daemon:
	@go build -o ${BINARY} ${PACKAGE_DAEMON}

run:
	@${BINARY} ${FLAGS}

test: test-unit test-postgres

test-unit:
	@go test -v ${PACKAGE_DAEMON}

test-postgres:
	@go test -v -tags=postgres ${PACKAGE_DAEMON} ${FLAGS}

get:
	@go get ${PACKAGE_DAEMON}