FROM golang:1.17-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache git

# placing files in GOROOT PATH
WORKDIR /usr/local/go/src/app/

# Pulling dependencies
COPY ./go.* ./
RUN go mod download

# Building stuff
COPY . ./
RUN go build -ldflags "-X main.Version=`git rev-parse HEAD`" -o app github.com/LasseJacobs/go-starter-kit/cmd

FROM alpine:3.7
RUN adduser -D -u 1000 runman

COPY --from=build /usr/local/go/src/app/app /usr/local/bin/app

EXPOSE 8080

USER runman
CMD ["app", "serve"]