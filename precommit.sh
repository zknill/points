#!/usr/bin/env bash

FILES=$(find . -name "*.go" | grep -v vendor/ | xargs -I % dirname % | sed 's/^.\///;s/[^.].*$/&\/*.go/;s/^\.$/*.go/' | sort -u)

echo "Running gofmt..."
res=$(gofmt -l ${FILES})
if [ -n "${res}" ]; then
    echo -e "format me please... \n${res}"
    exit 255
fi


echo "Running gometalinter..."
gometalinter --vendor -D gotype ./... --deadline=120s

echo "Running go test..."
go test -timeout 40s -race $(go list ./... | grep -v /vendor/)
