PKG=github.com/cyverse/irods-compare
VERSION=v0.1.0
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS?="-X '${PKG}/pkg/commons.clientVersion=${VERSION}' -X '${PKG}/pkg/commons.gitCommit=${GIT_COMMIT}' -X '${PKG}/pkg/commons.buildDate=${BUILD_DATE}'"
GO111MODULE=on
GOPROXY=direct
GOPATH=$(shell go env GOPATH)

.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags=${LDFLAGS} -o bin/irods-compare ./cmd/
	CGO_ENABLED=0 GOOS=windows go build -ldflags=${LDFLAGS} -o bin/irods-compare.exe ./cmd/
