PROTOC=/usr/local/bin/protoc
SERVICE=mnemosyne
PACKAGE=github.com/piotrkowalczuk/mnemosyne
PACKAGE_DAEMON=$(PACKAGE)/$(SERVICE)d
PACKAGE_TEST=$(PACKAGE)/$(SERVICE)test
BINARY=${SERVICE}d/${SERVICE}d

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
	-sp.connectionstring=$(MNEMOSYNE_STORAGE_POSTGRES_CONNECTION_STRING) \
	-sp.tablename=$(MNEMOSYNE_STORAGE_POSTGRES_TABLE_NAME)

CMD_TEST=go test -v -coverprofile=profile.out -covermode=atomic

.PHONY:	all proto build build-daemon run test test-unit test-postgres install package

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

rebuild: proto mocks build

run:
	@${BINARY} ${FLAGS}

test: test-unit test-postgres

test-unit:
	@${CMD_TEST} ${PACKAGE}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} -tags=unit ${PACKAGE_DAEMON}
	@cat profile.out >> coverage.txt && rm profile.out
	@${CMD_TEST} ${PACKAGE_TEST}

test-postgres:
	@${CMD_TEST} -tags=postgres ${PACKAGE_DAEMON} ${FLAGS}
	@cat profile.out >> coverage.txt && rm profile.out

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
		echo "error:\nPackaging requires effing package manager (fpm) to run.\nsee https://github.com/jordansissel/fpm\n"; \
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