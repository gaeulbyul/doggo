CLI_BIN := ./bin/doggo.bin
API_BIN := ./bin/doggo-api.bin

LAST_COMMIT := $(shell git rev-parse --short HEAD)
LAST_COMMIT_DATE := $(shell git show -s --format=%cI ${LAST_COMMIT})
VERSION := $(shell git describe --tags)
BUILDSTR := ${VERSION} | Commit ${LAST_COMMIT_DATE}-${LAST_COMMIT} | Build $(shell date --iso-8601=seconds)


.PHONY: build-cli
build-cli:
	go build -o ${CLI_BIN} -ldflags="-X 'main.buildString=${BUILDSTR}'" ./cmd/doggo/

.PHONY: build-api
build-api:
	go build -o ${API_BIN} -ldflags="-X 'main.buildString=${BUILDSTR}'" ./cmd/api/

.PHONY: build
build: build-api build-cli

.PHONY: run-cli
run-cli: build-cli ## Build and Execute the CLI binary after the build step.
	${CLI_BIN}

.PHONY: run-api
run-api: build-api ## Build and Execute the API binary after the build step.
	${API_BIN} --config config-api-sample.toml

.PHONY: clean
clean:
	go clean
	- rm -rf ./bin/

.PHONY: lint
lint:
	golangci-lint run
