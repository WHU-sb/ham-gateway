package gateway

//go:generate protoc --go_out=../gen --go_opt=paths=source_relative --go-grpc_out=../gen --go-grpc_opt=paths=source_relative --proto_path=.. ../proto/common.proto ../proto/course.proto ../proto/score.proto
