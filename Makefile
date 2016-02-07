PROTOC=/usr/local/bin/protoc
SERVICE=mnemosyne
PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test
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
          	    -sp.tablename=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE_NAME)

CMD_TEST=go test -v

.PHONY:	all proto build build-daemon run test test-unit test-postgres

all: proto build test run

proto:
	@${PROTOC} --proto_path=${GOPATH}/src \
	    --proto_path=. \
	    --proto_path=${GOPATH}/src/github.com/piotrkowalczuk/protot \
	    --go_out=Mprotot.proto=github.com/piotrkowalczuk/protot,plugins=grpc:. \
	    ${SERVICE}.proto
	@ls -al | grep "pb.go"

mocks:
	@mockery -all -output=${SERVICE}test -output_file=mocks.go -output_pkg_name=mnemosynetest

build: build-daemon

build-daemon:
	@go build -o ${BINARY} ${PACKAGE_DAEMON}

run:
	@${BINARY} ${FLAGS}

test: test-lib test-test test-daemon

test-lib:
	@${CMD_TEST} ${PACKAGE}

test-test:
	@${CMD_TEST} ${PACKAGE_TEST}

test-daemon:
	@${CMD_TEST} -tags=unit ${PACKAGE_DAEMON}
	@${CMD_TEST} -tags=postgres ${PACKAGE_DAEMON} -- ${FLAGS}

get:
	@go get github.com/stretchr/testify/...
	@go get github.com/onsi/ginkgo
	@go get github.com/onsi/gomega
	@go get ./...