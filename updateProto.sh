#!/bin/bash

rm -f internal/grpcs/proto/short.pb.go
rm -f internal/grpcs/proto/short_grpc.pb.go
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/grpcs/proto/short.proto
