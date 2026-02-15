FROM golang:1.25-alpine3.22 AS build_deps

RUN apk add git

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

FROM alpine:3.22

# Install minimal runtime
RUN apk add --no-cache ca-certificates

LABEL org.opencontainers.image.source="https://github.com/lukaarma/cert-manager-webhook-f5"

RUN adduser -D -u 1000 appuser
USER appuser

COPY --from=build /workspace/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook"]
