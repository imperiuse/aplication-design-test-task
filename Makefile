.EXPORT_ALL_VARIABLES:

-include .env

APP_NAME ?= aplication-design-test-task
APP_ENV ?= dev
APP_VERSION ?= dev

#BUILD_WITH_DEBUG    ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${NAME} -gcflags "all=-N -l" -ldflags '-v       -linkmode internal -extldflags \"-static\" -X ${GO_PACKAGE}/app/config.Version=${VERSION}' ${GO_PACKAGE}
#BUILD_WITHOUT_DEBUG ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/${NAME}
# We have to add an environment variable CGO_ENABLED=0 to disable dynamically links for a few dependencies. Normally we could not run Go applications from scratch because of this. We can also get rid of two more things in our binary. DWARF tables and annotations. The tables are needed for debuggers and the annotations for stack traces. Adding-ldflags="-s -w" removes them from our binary
BUILD_CMD ?= CGO_ENABLED=0 go build -tags=jsoniter -a -v -o bin/${APP_NAME} -ldflags '-v -w -s -linkmode auto -extldflags \"-static\" -X  main.AppName=${APP_NAME}  -X  main.AppVersion=${APP_VERSION}  -X  main.AppEnv=${APP_ENV}' ./cmd/${APP_NAME}
UPX_CMD ?= upx --best --lzma bin/${APP_NAME}

MACHINE_IP ?= 127.0.0.1
.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	@echo "Running build"
	${BUILD_CMD}
	@echo "Running UPX (zip) binary"
	${UPX_CMD}

.PHONY: lint
lint:
	@echo "Run golangci-lint"
	golangci-lint run -v ./...

.PHONY: tests
tests:
	@echo "Running go tests"
	go test -timeout 5m -v -race `go list ./... `

.PHONY: coverage
coverage:
	@echo "Running coverage.sh script"
	./tools/coverage.sh
