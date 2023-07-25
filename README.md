How to build the WASM
===

* `git submodule update --init`
* install docker
* `cd build && ./build.sh`` on a Unix-like environment (WSL might work too if it has access to the Docker)
  * note - it might take a bit of time on the CMake step
* in build/out/pdftops.wasm there is the compiled wasm file