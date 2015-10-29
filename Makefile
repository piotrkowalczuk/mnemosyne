PROTOC=/usr/local/bin/protoc
SERVICE=mnemosyne
PACKAGE=github.com/go-soa/mnemosyne
PACKAGE_SHARED=$(PACKAGE)/shared
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d

FLAGS=-h=$(MNEMOSYNE_HOST) \
      	    -p=$(MNEMOSYNE_PORT) \
      	    -s=$(MNEMOSYNE_SUBSYSTEM) \
      	    -lf=$(MNEMOSYNE_LOGGER_FORMAT) \
      	    -la=$(MNEMOSYNE_LOGGER_ADAPTER) \
      	    -ll=$(MNEMOSYNE_LOGGER_LEVEL) \
      	    -se=$(MNEMOSYNE_STORAGE_ENGINE) \
      	    -spcs=$(MNEMOSYNE_STORAGE_POSTGRES_CONNECTION_STRING) \
      	    -sptn=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE_NAME) \
      	    -spr=$(MNEMOSYNE_STORAGE_POSTGRES_RETRY)

.PHONY:	all proto build build-daemon run test test-unit test-postgres

all: proto build test run

proto:
	@${PROTOC} --proto_path=./shared --go_out=plugins=grpc:shared ./shared/${SERVICE}.proto
	@ls -al ./shared

build: build-daemon

build-daemon:
	@go build -o mnemosyned/mnemosyned ${PACKAGE_DAEMON}

run:
	@./mnemosyned/mnemosyned ${FLAGS}

test: test-unit test-postgres

test-unit:
	@go test -v ${PACKAGE_DAEMON} ${FLAGS}
	@go test -v ${PACKAGE_SHARED} ${FLAGS}

test-postgres:
	@go test -tags postgres -v ${PACKAGE_DAEMON} ${FLAGS}