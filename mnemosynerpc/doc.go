// Package mnemosynerpc ...
package mnemosynerpc

//go:generate protoc -I=. -I=../vendor --go_out=plugins=grpc:. session.proto
//go:generate goimports -w session.pb.go
//go:generate mockery -output=../mnemosynetest -output_file=mnemosynerpc.go -output_pkg_name=mnemosynetest -name=SessionManagerClient
//go:generate goimports -w ../mnemosynetest/mnemosynerpc.go
