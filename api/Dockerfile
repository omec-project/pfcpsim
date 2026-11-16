# SPDX-License-Identifier: Apache-2.0
# Copyright 2024-present Intel Corporation

FROM golang:1.27.4-bookworm@sha256:4fa7c968ba93f1635670522dbb01a2c78150443e1385746924566079f36570ea AS protobuf

RUN apt-get update && apt-get install -y --no-install-recommends unzip wget && \
    wget -q -O protoc-latest.zip https://github.com/protocolbuffers/protobuf/releases/download/v28.2/protoc-28.2-linux-x86_64.zip && \
    unzip protoc-latest.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf protoc3 protoc-latest.zip && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
    rm -rf /var/lib/apt/lists/*
