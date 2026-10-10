# syntax=docker/dockerfile:1.7
#
# Otter worker — background jobs (Sprint 19). Same base as the API (DD-203);
# go.sum-pinned dependencies and a self-probing HEALTHCHECK (DD-232).
# Build context: project root.

FROM golang:1.27-alpine AS deps

WORKDIR /src/server

COPY server/go.mod server/go.sum ./
RUN go mod download && go mod verify

FROM deps AS build

COPY server/ ./

ARG COMMIT_SHA=dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.CommitSHA=${COMMIT_SHA}" \
    -o /out/worker ./cmd/worker

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=build /out/worker /app

USER nonroot:nonroot

HEALTHCHECK --interval=30s --timeout=4s --retries=3 CMD ["/app", "healthcheck"]
ENTRYPOINT ["/app"]
