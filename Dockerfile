# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w \
      -X github.com/jjuanrivvera/alegra-cli/internal/version.Version=${VERSION} \
      -X github.com/jjuanrivvera/alegra-cli/internal/version.Commit=${COMMIT} \
      -X github.com/jjuanrivvera/alegra-cli/internal/version.BuildDate=${BUILD_DATE}" \
    -o /out/alegra ./cmd/alegra

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/alegra /usr/local/bin/alegra
ENTRYPOINT ["/usr/local/bin/alegra"]
