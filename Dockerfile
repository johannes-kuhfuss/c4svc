# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

# Setup ENV
WORKDIR /app

# Download prereqs
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy sources
COPY . .

# Build
RUN go build -o /c4svc

##
## Deploy
##
FROM alpine:latest

WORKDIR /

COPY --from=build /c4svc /c4svc

EXPOSE 8080

USER nobody:nogroup

ENTRYPOINT ["/c4svc"]
