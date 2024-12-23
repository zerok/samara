FROM golang:1.23.3-alpine AS builder
RUN mkdir -p /src
WORKDIR /src
COPY . /src

ARG VERSION

RUN find .
RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o samara -ldflags="-X 'main.version=${VERSION}'" ./cmd/samara

FROM alpine:3.20
COPY --from=builder /src/samara /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/samara"]
