// Package mnemosynerpc
package mnemosynerpc

//go:generate protoc -I=. -I=../vendor --go_out=plugins=grpc:. mnemosyne.proto
//go:generate goimports -w mnemosyne.pb.go
//go:generate mockery -output=../mnemosynetest -output_file=mnemosynerpc.go -output_pkg_name=mnemosynetest -name=SessionManagerClient
