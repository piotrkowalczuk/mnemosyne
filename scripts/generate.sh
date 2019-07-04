SERVICE=mnemosyne
PROTO_INCLUDE="-I=/usr/include -I=."

goimports -w ./${SERVICE}rpc
mockery -case=underscore -dir=./${SERVICE}rpc -all -output=./${SERVICE}test -outpkg=${SERVICE}test
mockery -case=underscore -dir=./${SERVICE}d -all -case=underscore -output=./${SERVICE}d
mockery -case=underscore -dir=./internal/storage -all -output=./internal/storage/storagemock -outpkg=storagemock
gofmt -w -r 'google_protobuf1.Empty -> empty.Empty' ./${SERVICE}test