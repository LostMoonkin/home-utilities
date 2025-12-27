# Stage1 build
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git bash

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a static binary that works on Alpine
RUN /bin/bash ./build.sh

# Stage2 deploy
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /homeserver

COPY --from=builder /app/build/homeserver .

EXPOSE 8786

CMD ["./homeserver"]
