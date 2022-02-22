package helloworld

//go:generate go install ../../../cmd/protoc-gen-gogo-nrpc
//go:generate go install github.com/gogo/protobuf/protoc-gen-gogo
//go:generate protoc -I. -I../../.. -I../../../third_party -I ../../../third_party/gogoproto --gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. --gogo_opt=paths=source_relative --gogo-nrpc_out . helloworld.proto
