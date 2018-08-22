FROM golang:alpine AS build-stage

ENV GOROOT=/usr/local/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PROJPATH=/gopath/src/github.com/gzsunrun/ansible-manager/

RUN apk add -U -q --no-progress build-base

ADD . /gopath/src/github.com/gzsunrun/ansible-manager
WORKDIR /gopath/src/github.com/gzsunrun/ansible-manager

RUN make build



FROM alpine:latest

RUN apk add -U -q --no-progress ansible

COPY --from=build-stage /gopath/src/github.com/gzsunrun/ansible-manager/ansible-manager /usr/local/bin/

ENV PATH=$PATH:/usr/local/bin

WORKDIR /usr/local/bin

CMD ["./ansible-manager"]
