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
rm -f gin-bin || true
gin --all run main.go
