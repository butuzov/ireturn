# --- Required ----------------------------------------------------------------
export PATH   := $(PWD)/bin:$(PATH)                    # ./bin to $PATH
export SHELL  := bash                                  # Default Shell

build:
	@ go build -trimpath -o bin/ireturn ./cmd/ireturn/

tests:
	go test -v -count=1 -race \
		-failfast \
		-parallel=2 \
		-timeout=1m \
		-covermode=atomic \
		-coverprofile=coverage.cov ./...

lints:
	golangci-lint run --no-config ./... -D deadcode

cover:
	go tool cover -html=coverage.cov

install:
	go install -trimpath -v -ldflags="-w -s" ./cmd/ireturn/

bin/goreleaser:
	@curl -Ls https://github.com/goreleaser/goreleaser/releases/download/v1.18.2/goreleaser_Darwin_all.tar.gz | tar -zOxf - goreleaser > ./bin/goreleaser
	chmod 0755 ./bin/goreleaser

test-release: bin/goreleaser
	goreleaser release --help
	goreleaser release -f .goreleaser.yaml \
		--skip-validate --skip-publish --clean
