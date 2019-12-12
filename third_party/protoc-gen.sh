#!/bin/bash
protoc -I third_party \
    -I api/proto/v0 \
    --go_out=plugins=grpc:pkg/api/v0 \
    account-service.proto
