#!/bin/bash
go get ./...
rm -rf ./build/jp-*
for goos in darwin linux windows freebsd; do
	export GOOS="$goos"
	for goarch in 386 amd64; do
		export GOARCH="$goarch"
		go build -v -o build/jp-$GOOS-$GOARCH
	done
done
