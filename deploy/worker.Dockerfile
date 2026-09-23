# syntax=docker/dockerfile:1.7
#
# Otter Worker — background job processor.
# Same multi-stage + distroless pattern as API.

# ---- dependencies ----
FROM golang:1.25-alpine AS deps

RUN apk add --no-cache git

WORKDIR /src/server

COPY server/go.mod ./

RUN go mod tidy && go mod download && go mod verify

# ---- build ----
FROM deps AS build

COPY server/ ./

ARG COMMIT_SHA=dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/worker ./cmd/worker

# ---- runtime ----
FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /out/worker /app

USER nonroot:nonroot

ENTRYPOINT ["/app"]
