#!/bin/bash

# fail on error
set -e

# =============================================================================================
if [[ "$(basename $PWD)" == "scripts" ]]; then
    cd ..
fi
echo $PWD

# =============================================================================================
source .env
source ~/.config/ir.conf || true

# =============================================================================================
echo "developing iRcollector ..."
killall gin-bin || true
killall iRcollector || true
rm -f gin-bin || true
#gin --all run main.go

rm -f iRcollector || true
GOARCH=amd64 GOOS=linux go build -i -o iRcollector
./iRcollector
