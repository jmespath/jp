#!/bin/bash

set -e

cd build
shasum -a 256 ./jp-* > jp-checksums.sha256
gpg --clearsign --output jp-checksums.sha256.asc jp-checksums.sha256
rm jp-checksums.sha256
