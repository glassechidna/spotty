#!/bin/bash
set -euxo pipefail

go get github.com/mitchellh/gox

curl -z dl/upx.txz -o dl/upx.txz -L https://github.com/upx/upx/releases/download/v3.94/upx-3.94-amd64_linux.tar.xz
tar -xvf dl/upx.txz

curl -z dl/glide.tgz -o dl/glide.tgz -L https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
tar -xvf dl/glide.tgz
