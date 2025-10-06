# SPDX-License-Identifier: AGPL-3.0-or-later
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git make gcc musl-dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/fractal /usr/local/bin/fractal
COPY fractal.toml.example /etc/fractal.toml.example
EXPOSE 8899
ENTRYPOINT ["fractal"]
