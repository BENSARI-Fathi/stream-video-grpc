#!/bin/bash

protoc --go-grpc_out=. --go_out=. streamVod/streamVodpb/streamVod.proto
