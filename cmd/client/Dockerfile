FROM golang:1.24.1-alpine3.20 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download
COPY ./internal ./internal
RUN CGO_ENABLED=0 go test -short ./internal...
COPY ./cmd/client ./cmd/client
RUN CGO_ENABLED=0 go test -short ./cmd/client/...
RUN CGO_ENABLED=0 go build -o client ./cmd/client


FROM alpine:3.20 AS final
WORKDIR /app
COPY --from=build /build/client ./
ENTRYPOINT ["/app/client"]
