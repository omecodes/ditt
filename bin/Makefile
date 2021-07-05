VERSION=0.0.1-dev
APP_NAME=ditt-api-server
TODAY:=$(shell date -u +%Y-%m-%dT%H:%M:%S)
TIMESTAMP:=$(shell date -u +%Y%m%d%H%M%S)
GITREV:=$(shell git rev-parse HEAD)

LINK_FLAGS="-X github.com/omecodes/ditt/info.Version=${VERSION} -X github.com/omecodes/ditt/info.ApplicationName=${APP_NAME} -X github.com/omecodes/ditt/info.BuildStamp=${TODAY} -X github.com/omecodes/ditt/info.BuildRevision=${GITREV}"

all: tests build

tests:
	cd .. && go test ./...

build:
	go build --tags json1 -o ${APP_NAME} -ldflags=${LINK_FLAGS} api-server.go

