# --- Required ----------------------------------------------------------------
export PATH   := $(PWD)/bin:$(PATH)                    # ./bin to $PATH
export SHELL  := bash                                  # Default Shell

GOPKGS := $(shell go list ./... | grep -vE "(testdata)" | tr -s '\n' ',' | sed 's/.\{1\}$$//' )


build:
	@ go build -trimpath -o bin/ireturn ./cmd/ireturn/

tests:
	go test -v -count=1 -race -covermode=atomic \
		-coverpkg=$(GOPKGS) -coverprofile=coverage.cov ./...
