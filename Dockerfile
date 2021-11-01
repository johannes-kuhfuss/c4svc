# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy all files & directories
COPY app/ ./
COPY config/ ./
COPY controllers ./
COPY domain/ ./
COPY providers/ ./
COPY services/ ./
COPY utils/ ./
COPY *.go ./

RUN go build -o /c4svc

##
## Deploy
##
FROM gcr.io/distroless/base

WORKDIR /

COPY --from=build /c4svc /c4svc

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/c4svc"]