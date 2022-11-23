########## Builder ##########
FROM golang:1.19-alpine AS builder

# Install the latest version of Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy local source
COPY . /build
WORKDIR /build

# Build the binary
RUN go build ./cmd/example/

########## Runtime ##########
FROM alpine:3

WORKDIR /

# Copy binaries from builder
COPY --from=builder /build /
COPY --from=builder /go/bin/dlv /

# Expose both the apps port and debugger
EXPOSE 8080 40000

# Start Delve
CMD ./dlv --listen=:40000 --headless=true --api-version=2 --log exec ./example