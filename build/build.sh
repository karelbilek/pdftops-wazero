#!/usr/bin/env bash

docker build -t pdftops-wazero-build .
docker run -v $(pwd)/out:/xpdf/build/out:z pdftops-wazero-build /usr/bin/cp /xpdf/build/xpdf/pdftops.wasm /xpdf/build/out/pdftops.wasm