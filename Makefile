BIN_SERVER=mantrae
BIN_AGENT=mantrae_agent

VERSION=$(shell git describe --tags --abbrev=0)
DATE=$(shell date -u +%Y-%m-%d)
COMMIT=$(shell git rev-parse --short HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X github.com/MizuchiLabs/mantrae/internal/util.Version=${VERSION} -X github.com/MizuchiLabs/mantrae/internal/util.BuildDate=${DATE} -X github.com/MizuchiLabs/mantrae/internal/util.Commit=${COMMIT}"

all: clean build

.PHONY: clean
clean:
	rm -rf $(PWD)/$(BIN) $(PWD)/$(BIN)-agent $(PWD)/web/build $(PWD)/builds

.PHONY: audit
audit-security:
	- gosec --exclude=G104 ./...
	- govulncheck -show=color ./...
	- staticcheck -checks=all -f=stylish ./...

audit: audit-security
	go fmt ./...
	go vet ./...
	go mod tidy
	go mod verify
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

.PHONY: build
build-web: $(PWD)/web/build
$(PWD)/web/build: $(shell find web/src -type f) web/package.json web/pnpm-lock.yaml
	cd web && pnpm install && pnpm run build
	touch $(PWD)/web/build

build-server: build-web
	 go build $(LDFLAGS) -o $(BIN_SERVER) main.go

build-agent:
	 go build $(LDFLAGS) -o $(BIN_AGENT) agent/cmd/main.go

docker-server: snapshot
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=${VERSION} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg DATE=${DATE} \
		-t ghcr.io/mizuchilabs/mantrae:latest \
		.

docker-agent: snapshot
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=${VERSION} \
		--build-arg COMMIT=${COMMIT} \
		--build-arg DATE=${DATE} \
		-t ghcr.io/mizuchilabs/mantrae-agent:latest \
		.

docker-push:
	docker push ghcr.io/mizuchilabs/mantrae:latest

.PHONY: release
release:
	goreleaser release --clean --skip=validate

.PHONY: snapshot
snapshot:
	goreleaser release --clean --snapshot

.PHONY: upgrade
upgrade:
	go get -u && go mod tidy
	cd web && pnpm update

.PHONY: db-up
db-up:
	goose sqlite3 mantrae.db -dir internal/db/migrations up

.PHONY: db-down
db-down:
	goose sqlite3 mantrae.db -dir internal/db/migrations down

.PHONY: db-reset
db-reset:
	rm -f mantrae.db
	goose sqlite3 mantrae.db -dir internal/db/migrations up

.PHONY: db-status
db-status:
	goose sqlite3 mantrae.db -dir internal/db/migrations status

.PHONY: run
run-server:
	go run main.go

run-web:
	cd web && npm run dev
