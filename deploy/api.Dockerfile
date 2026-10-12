# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine AS build

WORKDIR /src/server

# Download Go dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy Go source
COPY server/ ./

ARG COMMIT_SHA=dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/api ./cmd/api

# Runtime
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=build /out/api /app/api
COPY content/ /app/content/

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/api"]
