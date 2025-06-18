BIN_SERVER=mantrae
BIN_AGENT=mae

VERSION=$(shell git describe --tags --abbrev=0)
DATE=$(shell date -u +%Y-%m-%d)
COMMIT=$(shell git rev-parse --short HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X github.com/MizuchiLabs/mantrae/pkg/build.Version=${VERSION} -X github.com/MizuchiLabs/mantrae/pkg/build.Date=${DATE} -X github.com/MizuchiLabs/mantrae/pkg/build.Commit=${COMMIT}"

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

build-server:
	go generate ./...
	go build $(LDFLAGS) -o $(BIN_SERVER) main.go

build-agent:
	go build $(LDFLAGS) -o $(BIN_AGENT) agent/cmd/main.go

docker-local:
	go generate ./...
	KO_DOCKER_REPO=ko.local/mantrae ko build . --bare
	KO_DOCKER_REPO=ko.local/mantrae-agent ko build ./agent/cmd --bare

docker-release:
	go generate ./...
	KO_DOCKER_REPO=ghcr.io/mizuchilabs/mantrae ko build . --bare
	KO_DOCKER_REPO=ghcr.io/mizuchilabs/mantrae-agent ko build ./agent/cmd --bare

.PHONY: release
release:
	go generate ./...
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
	goose sqlite3 data/mantrae.db -dir internal/store/migrations up-by-one

.PHONY: db-down
db-down:
	goose sqlite3 data/mantrae.db -dir internal/store/migrations down

.PHONY: db-reset
db-reset:
	rm -f mantrae.db
	goose sqlite3 data/mantrae.db -dir internal/store/migrations up

.PHONY: db-status
db-status:
	goose sqlite3 data/mantrae.db -dir internal/store/migrations status

.PHONY: run
run-server:
	go run main.go

run-web:
	cd web && npm run dev
