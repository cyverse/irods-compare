PKG=github.com/cyverse/irods-compare
VERSION=v0.1.1
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
	CGO_ENABLED=0 go build -ldflags=${LDFLAGS} -o bin/irods-compare ./cmd/

.PHONY: build-release
build-release:
	mkdir -p release

	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags=${LDFLAGS} -o release/irods-compare ./cmd/
	cd release && tar cvf irods_compare_i386_linux_${VERSION}.tar irods-compare && cd ..
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=${LDFLAGS} -o release/irods-compare ./cmd/
	cd release && tar cvf irods_compare_amd64_linux_${VERSION}.tar irods-compare && cd ..
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags=${LDFLAGS} -o release/irods-compare ./cmd/
	cd release && tar cvf irods_compare_arm_linux_${VERSION}.tar irods-compare && cd ..
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags=${LDFLAGS} -o release/irods-compare ./cmd/
	cd release && tar cvf irods_compare_arm64_linux_${VERSION}.tar irods-compare && cd ..

	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags=${LDFLAGS} -o release/irods-compare ./cmd/
	cd release && tar cvf irods_compare_amd64_darwin_${VERSION}.tar irods-compare && cd ..
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags=${LDFLAGS} -o release/irods-compare ./cmd/
	cd release && tar cvf irods_compare_arm64_darwin_${VERSION}.tar irods-compare && cd ..

	rm release/irods-compare

	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags=${LDFLAGS} -o release/irods-compare.exe ./cmd/
	cd release && tar cvf irods_compare_i386_windows_${VERSION}.tar irods-compare.exe && cd ..
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags=${LDFLAGS} -o release/irods-compare.exe ./cmd/
	cd release && tar cvf irods_compare_amd64_windows_${VERSION}.tar irods-compare.exe && cd ..
	rm release/irods-compare.exe