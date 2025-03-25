FROM golang:1.24.1-alpine3.20 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download
COPY ./internal ./internal
RUN CGO_ENABLED=0 go test -short ./internal...
COPY ./cmd/server ./cmd/server
RUN CGO_ENABLED=0 go test -short ./cmd/server/...
RUN CGO_ENABLED=0 go build -o server ./cmd/server


FROM alpine:3.20 AS final
WORKDIR /app
COPY quotes.yaml ./
COPY --from=build /build/server ./
ENTRYPOINT ["/app/server"]
EXPOSE 1337/tcp
