#!/bin/bash
go get ./...
rm -rf ./build/jp-*
for goos in darwin linux windows freebsd; do
	export GOOS="$goos"
	for goarch in 386 amd64; do
		export GOARCH="$goarch"
		echo "Building for $GOOS/$GOARCH"
		go build -v -o build/jp-$GOOS-$GOARCH 2>/dev/null
	done
done
# Also build for ARM7/linux
export GOOS=linux
export GOARCH=arm
export GOARM=7
echo "Building for $GOOS/$GOARCH/$GOARM"
go build -v -o build/jp-$GOOS-$GOARCH-arm$GOARM 2> /dev/null
