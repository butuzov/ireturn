FROM golang:1.26-alpine AS builder

WORKDIR /build
RUN   apk add --no-cache upx
COPY  go.mod  .
RUN   go mod download -x
COPY  . .
RUN   go build -trimpath -o bin/ireturn ./cmd/ireturn
RUN   upx --brute /build/bin/ireturn


FROM alpine:latest
WORKDIR    /
COPY       --from=builder /build/bin/ireturn ireturn
VOLUME     /app
WORKDIR    /app
ENTRYPOINT ["/ireturn" ]
