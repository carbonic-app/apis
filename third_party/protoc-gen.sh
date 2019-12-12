#!/bin/bash
protoc -I api/proto/v0 \
    --go_out=plugins=grpc:pkg/api/v0 \
    token.proto \
    account_service.proto
