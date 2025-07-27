Pstopdf tool from xpdfreader, converted to wasm, run with WASM+wazero+fake memory FS.

It's one part of PDF-to-PDF/A conversion. (We use this on some weird PDFs that refuse to be 
converted just with ghostscript; most are happy with ghostscript, some are not.)

We are actually building all xpdfreader tools to wasm; but we then use just pdftops.

How to use
===

See example_test.go
```
package main

import (
    "context"
    "fmt"

    pdftops "github.com/karelbilek/pdftops-wazero"
)

func main() {
	g, err = pdftops.New(context.Background())
	if err != nil {
		panic(err)
	}

	in, err := os.ReadFile("some.pdf")
	if err != nil {
		panic(err)
	}

	out, err := g.ConvertPDFToPS(context.Background(), in)
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