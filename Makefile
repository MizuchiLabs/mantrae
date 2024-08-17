BIN=mantrae

VERSION=$(shell git describe --tags)
DATE=$(shell date -u +%Y-%m-%d)
COMMIT=$(shell git rev-parse --short HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X github.com/MizuchiLabs/mantrae/tools/updater.Version=${VERSION} -X github.com/MizuchiLabs/mantrae/tools/updater.BuildDate=${DATE} -X github.com/MizuchiLabs/mantrae/tools/updater.Commit=${COMMIT}"

all: clean build

.PHONY: clean
clean:
	rm -rf $(PWD)/bin/ $(PWD)/web/build

.PHONY: audit
audit:
	go fmt ./...
	go vet ./...
	go mod tidy
	go mod verify
	- gosec --exclude=G104 ./...
	- govulncheck -show=color ./...
	- staticcheck -checks=all -f=stylish ./...

.PHONY: build
build: audit
	cd web && pnpm install && pnpm run build
	go build $(LDFLAGS) -o bin/$(BIN) main.go && upx bin/$(BIN)

build-fast:
	go build $(LDFLAGS) -o bin/$(BIN) main.go

.PHONY: docker
docker:
	docker build \
		--label "org.opencontainers.image.source=https://github.com/MizuchiLabs/mantrae" \
		--label "org.opencontainers.image.description=Mantrae" \
		--label "org.opencontainers.image.version=${VERSION}" \
		--label "org.opencontainers.image.revision=${COMMIT}" \
		--label "org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		--label "org.opencontainers.image.licenses=Apache-2.0" \
		-t ghcr.io/mizuchilabs/mantrae:${VERSION} .
	docker tag ghcr.io/mizuchilabs/mantrae:${VERSION} ghcr.io/mizuchilabs/mantrae:latest

.PHONY: upgrade
upgrade:
	go get -u && go mod tidy
	cd web && pnpm update

.PHONY: run
run-server:
	go run main.go

run-web:
	cd web && npm run dev
