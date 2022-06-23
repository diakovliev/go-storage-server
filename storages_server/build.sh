#! /bin/sh

repo_root="$(git rev-parse --show-toplevel)"

export GO111MODULE=off
export GOPATH=${repo_root}/.gopath
export GOBIN=${GOPATH}/bin

mkdir -p ${GOPATH}/bin

go get
go build
