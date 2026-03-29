# Build stage
FROM golang:1.25-alpine AS build
RUN apk add --no-cache build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/sanitize ./cmd/sanitize

# Runtime stage
FROM alpine:3.21

LABEL org.opencontainers.image.title="sanitize"
LABEL org.opencontainers.image.description="Structured-text linter and heuristic fixer with web playground"
LABEL org.opencontainers.image.source="https://github.com/go-go-golems/sanitize"

RUN apk add --no-cache ca-certificates
RUN adduser -D -H -u 10001 appuser
USER 10001
WORKDIR /
COPY --from=build /out/sanitize /sanitize
EXPOSE 8080
ENTRYPOINT ["/sanitize", "serve"]
