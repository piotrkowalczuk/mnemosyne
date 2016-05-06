PROTOC=/usr/local/bin/protoc
SERVICE=mnemosyne
PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
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
	-subsystem=$(MNEMOSYNE_SUBSYSTEM) \
	-namespace=$(MNEMOSYNE_NAMESPACE) \
	-l.format=$(MNEMOSYNE_LOGGER_FORMAT) \
	-l.adapter=$(MNEMOSYNE_LOGGER_ADAPTER) \
	-l.level=$(MNEMOSYNE_LOGGER_LEVEL) \
	-m.engine=$(MNEMOSYNE_MONITORING_ENGINE) \
	-s.engine=$(MNEMOSYNE_STORAGE_ENGINE) \
	-s.p.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS) \
	-s.p.table=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE)

CMD_TEST=go test -coverprofile=profile.out -covermode=atomic

.PHONY:	all proto build build-daemon run test test-short install package

all: proto build test run

proto:
	@${PROTOC} --proto_path=${GOPATH}/src \
	    --proto_path=. \
	    --go_out=plugins=grpc:. \
	    ${SERVICE}.proto
	@ls -al | grep "pb.go"

mocks:
	@mockery -all -output=${SERVICE}test -output_file=mocks.go -output_pkg_name=mnemosynetest

build: build-daemon

build-daemon:
	@go build -o .tmp/${SERVICE}d ${PACKAGE_CMD_DAEMON}

rebuild: proto mocks build

run:
	@.tmp/${SERVICE}d ${FLAGS}

test-short:
	@${CMD_TEST} -short ${PACKAGE}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} -short ${PACKAGE_DAEMON}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} -short ${PACKAGE_TEST}

test:
	@${CMD_TEST} ${PACKAGE}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} ${PACKAGE_DAEMON} -s.p.address=$(MNEMOSYNE_STORAGE_POSTGRES_ADDRESS)
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} ${PACKAGE_TEST}


get:
	@go get github.com/smartystreets/goconvey/...
	@go get ./...

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