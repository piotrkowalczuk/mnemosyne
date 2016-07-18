PROTOC=/usr/local/bin/protoc
SERVICE=mnemosyne
PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
PACKAGE_RPC=$(PACKAGE)/$(SERVICE)rpc
PACKAGE_CMD_DAEMON=$(PACKAGE)/cmd/$(SERVICE)d
BINARY_CMD_DAEMON=cmd/${SERVICE}d/${SERVICE}d

#packaging
DIST_PACKAGE_BUILD_DIR=temp
DIST_PACKAGE_DIR=dist
DIST_PACKAGE_TYPE=deb
DIST_PREFIX=/usr
DIST_BINDIR=${DESTDIR}${DIST_PREFIX}/bin

FLAGS=-host=$(MNEMOSYNE_HOST) \
	-port=$(MNEMOSYNE_PORT) \
	-ttl=$(MNEMOSYNE_TTL) \
	-ttc=$(MNEMOSYNE_TTC) \
	-subsystem=$(MNEMOSYNE_SUBSYSTEM) \
	-namespace=$(MNEMOSYNE_NAMESPACE) \
	-log.format=$(MNEMOSYNE_LOGGER_FORMAT) \
	-log.adapter=$(MNEMOSYNE_LOGGER_ADAPTER) \
	-log.level=$(MNEMOSYNE_LOGGER_LEVEL) \
	-monitoring=$(MNEMOSYNE_MONITORING_ENGINE) \
	-storage.engine=$(MNEMOSYNE_STORAGE_ENGINE) \
	-storage.postgres.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS) \
	-storage.postgres.table=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE)

CMD_TEST=go test -race -coverprofile=.tmp/profile.out -covermode=atomic

.PHONY:	all gen build rebuild run test test-short get install

all: get install

gen:
	@go generate .
	@go generate ./${SERVICE}d
	@go generate ./${SERVICE}rpc
	@ls -al ./${SERVICE}rpc | grep "pb.go"

build:
	@go build -o .tmp/${SERVICE}d ${PACKAGE_CMD_DAEMON}

rebuild: gen build

run:
	@.tmp/${SERVICE}d ${FLAGS}

test-short:
	@${CMD_TEST} -short ${PACKAGE}
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} -short ${PACKAGE_DAEMON}
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} -short ${PACKAGE_RPC}
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} -short ${PACKAGE_TEST}

test:
	@${CMD_TEST} ${PACKAGE} -storage.postgres.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS)
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} ${PACKAGE_DAEMON} -storage.postgres.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS)
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} ${PACKAGE_RPC}
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} ${PACKAGE_TEST}

get:
	@go get github.com/Masterminds/glide
	@go get github.com/smartystreets/goconvey/...
	@glide --no-color install

install:
	@go install ${PACKAGE_CMD_DAEMON}