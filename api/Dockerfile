# SPDX-License-Identifier: Apache-2.0
# Copyright 2024-present Intel Corporation

FROM golang:1.27.2-bookworm@sha256:5cf287a799e6b94384bad13d16b14904c531f51ba65792237e122ce42b392f61 AS protobuf

RUN apt-get update && apt-get install -y unzip wget && \
    wget -O protoc-latest.zip https://github.com/protocolbuffers/protobuf/releases/download/v28.2/protoc-28.2-linux-x86_64.zip && \
    unzip protoc-latest.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf protoc3 protoc-latest.zip && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
