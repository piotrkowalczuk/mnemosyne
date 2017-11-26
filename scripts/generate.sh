#!/bin/sh

SERVICE=mnemosyne
PROTO_INCLUDE="-I=/usr/include -I=."

protoc ${PROTO_INCLUDE} --go_out=plugins=grpc:${GOPATH}/src ${SERVICE}rpc/*.proto
python -m grpc_tools.protoc ${PROTO_INCLUDE} --python_out=. --grpc_python_out=. ${SERVICE}rpc/*.proto
goimports -w ./${SERVICE}rpc
mockery -case=underscore -dir=./${SERVICE}rpc -all -output=./${SERVICE}test -outpkg=${SERVICE}test
mockery -case=underscore -dir=./${SERVICE}d -all -inpkg
mockery -case=underscore -dir=./internal/storage -all -output=./internal/storage/storagemock -outpkg=storagemock
gofmt -w -r 'google_protobuf1.Empty -> empty.Empty' ./${SERVICE}test


#mockery -dir=./internal/model -name=.*Provider -output=./internal/model -inpkg
#mockery -dir=./internal/model -name=Rows -output=./internal/model -inpkg
#mockery -dir=./internal/event -name=.*Dispatcher -output=./internal/event -inpkg
#mockery -dir=./creativeservd -name=envelope -output=./creativeservd -inpkg
#mockery -dir=./creativeservd -name=DimensionLocker -output=./creativeservd -inpkg
#goimports -w ./${SERVICE}rpc
#goimports -w ./${SERVICE}test
#ls -lha ./${SERVICE}rpc | grep pb.go
#ls -lha ./${SERVICE}rpc | grep pb2.py
#ls -lha ./${SERVICE}rpc | grep grpc.py


#//go:generate protoc -I=. -I=/usr/include -I=../vendor --go_out=plugins=grpc:${GOPATH}/src *.proto
#//go:generate python -m grpc_tools.protoc -I=. -I=/usr/include -I=../vendor --python_out=. --grpc_python_out=${GOPATH}/src *.proto
#//go:generate mockery -output=../mnemosynetest -outpkg=mnemosynetest -all -case=underscore
#//go:generate goimports -w ../mnemosynetest