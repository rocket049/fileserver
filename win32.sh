#!/usr/bin/env bash
export CGO_ENABLED=0
export GOOS=windows
export GOARCH=386
go build -ldflags -s
