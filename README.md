Convert PDF to PS in cgo-less portable go, just with WASM+wazero.

It's one part of PDF-to-PDF/A conversion.

Note that the code writes to filesystem for temporary files; the local filesystem should not be full or read-only.

How to use
===

```
package main

import (
    pdftops "github.com/karelbilek/pdftops-wazero"
)

func main() {
	in, err := os.ReadFile("some.pdf")
	if err != nil {
		panic(err)
	}

	out, err := pdftops.ConvertPDFToPS(context.Background(), in)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("someout.ps", out, 0o666)
	if err != nil {
		panic(err)
	}
}


```

Copyright
===
(C) 2023 Karel Bilek, Jeroen Bobbeldijk (https://github.com/jerbob92)

xpdfreader (C) 1996-2022 Glyph & Cog, LLC. (from https://www.xpdfreader.com/download.html )

ghostscript fonts Copyright (c) 2001- Valek Filippov (from https://packages.ubuntu.com/focal/gsfonts)

The whole code is GPLv2.

How to build the WASM
===

* `git submodule update --init`
* install docker
* `cd build && ./build.sh` on a Unix-like environment (WSL might work too if it has access to the Docker)
  * note - it might take a bit of time on the CMake step
* in build/out/pdftops.wasm there is the compiled wasm file