# syntax=docker/dockerfile:1

##
## Build
##
###FROM golang:1.17-alpine AS build

# Setup ENV
###ENV GOPATH=/app
###ENV REPO_URL=github.com/johannes-kuhfuss/c4svc
###ENV APP_PATH=${GOPATH}/src/${REPO_URL}
###ENV WORKPATH=${APP_PATH}/src

# Copy
###COPY src ${WORKPATH}
###WORKDIR ${WORKPATH}

# Build
####RUN go build -o c4svc .

##
## Deploy
##
FROM gcr.io/distroless/base

WORKDIR /

COPY c4svc /c4svc

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/c4svc"]
