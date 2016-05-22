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
	-l.format=$(MNEMOSYNE_LOGGER_FORMAT) \
	-l.adapter=$(MNEMOSYNE_LOGGER_ADAPTER) \
	-l.level=$(MNEMOSYNE_LOGGER_LEVEL) \
	-m.engine=$(MNEMOSYNE_MONITORING_ENGINE) \
	-s.engine=$(MNEMOSYNE_STORAGE_ENGINE) \
	-s.p.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS) \
	-s.p.table=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE)

CMD_TEST=go test -race -coverprofile=.tmp/profile.out -covermode=atomic

.PHONY:	all proto build rebuild mocks run test test-short get install package

all: proto build test run

proto:
	@${PROTOC} --proto_path=${GOPATH}/src \
	    --proto_path=. \
	    --go_out=plugins=grpc:. \
	    ${SERVICE}rpc/${SERVICE}.proto
	@ls -al ./${SERVICE}rpc | grep "pb.go"

mocks:
	@mockery -all -output=${SERVICE}test -output_file=mocks.go -output_pkg_name=mnemosynetest

build:
	@go build -o .tmp/${SERVICE}d ${PACKAGE_CMD_DAEMON}

rebuild: proto mocks build

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
	@${CMD_TEST} ${PACKAGE} -s.p.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS)
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} ${PACKAGE_DAEMON} -s.p.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS)
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} ${PACKAGE_RPC}
	@cat .tmp/profile.out >> .tmp/coverage.txt && rm .tmp/profile.out
	@${CMD_TEST} ${PACKAGE_TEST}

get:
	@go get github.com/Masterminds/glide
	@go get github.com/smartystreets/goconvey/...
	@glide install

install: build
	#install binary
	install -Dm 755 ${BINARY} ${DIST_BINDIR}/${SERVICE}
	#install config file
	install -Dm 644 scripts/${SERVICE}.env ${DESTDIR}/etc/${SERVICE}.env
	#install init script
	install -Dm 755 scripts/${SERVICE}.service ${DESTDIR}/etc/systemd/system/${SERVICE}.service

package:
	# export DIST_PACKAGE_TYPE to vary package type (e.g. deb, tar, rpm)
	@if [ -z "$(shell which fpm 2>/dev/null)" ]; then \
		echo "error:\nPackagings requires effing package manager (fpm) to run.\nsee https://github.com/jordansissel/fpm\n"; \
		exit 1; \
	fi

	#run make install against the packaging dir
	mkdir -p ${DIST_PACKAGE_BUILD_DIR} && $(MAKE) install DESTDIR=${DIST_PACKAGE_BUILD_DIR}

	#clean
	mkdir -p ${DIST_PACKAGE_DIR} && rm -f ${DIST_PACKAGE_DIR}/*.${DIST_PACKAGE_TYPE}

	#build package
	fpm --rpm-os linux \
		-s dir \
		-p dist \
		-t ${DIST_PACKAGE_TYPE} \
		-n ${SERVICE} \
		-v `${DIST_PACKAGE_BUILD_DIR}${DIST_PREFIX}/bin/${SERVICE} -version` \
		--config-files /etc/${SERVICE}.env \
		--config-files /etc/systemd/system/${SERVICE}.service \
		-C ${DIST_PACKAGE_BUILD_DIR} .