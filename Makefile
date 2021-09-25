# --- Required ----------------------------------------------------------------
export PATH   := $(PWD)/bin:$(PATH)                    # ./bin to $PATH
export SHELL  := bash                                  # Default Shell

GOPKGS := $(shell go list ./... | grep -vE "(cmd|testdata)" | tr -s '\n' ',' | sed 's/.\{1\}$$//' )


build:
	@ go build -trimpath -o bin/ireturn ./cmd/ireturn/

tests:
	go test -v -count=1 -race \
		-failfast \
		-parallel=2 \
		-timeout=1m \
		-covermode=atomic \
		-coverpkg=$(GOPKGS) -coverprofile=coverage.cov ./...

lints:
	golangci-lint run --no-config ./... -D deadcode

cover:
	go tool cover -html=coverage.cov

install:
	go install -trimpath -v -ldflags="-w -s" ./cmd/ireturn/

bin/goreleaser:
	curl -fsSL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh |  sh -s -- -b bin

test-release: bin/goreleaser
	goreleaser release --help
	goreleaser release -f .goreleaser.yaml \
		--skip-validate --skip-publish --rm-dist
