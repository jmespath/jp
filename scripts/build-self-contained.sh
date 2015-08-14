#!/bin/bash
# This is a self contained shell script for building jp without having to set
# up GOPATH.  You just need go installed.
tempdir="$(mktemp -d -t jpbuild)"
tempgopath="$tempdir/go"
jppath="${tempgopath}/src/github.com/jmespath"
fakerepo="$jppath/jp"
mkdir -p $jppath
ln -s "$(pwd)" "$jppath/jp"
export GOPATH="$tempgopath"
cd "$fakerepo"
go build
