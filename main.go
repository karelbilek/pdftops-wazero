package pdftopswazero

import (
	"context"
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"

	"strings"

	"github.com/karelbilek/wazero-fs-tools/memfs"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/experimental/sysfs"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed build/out/pdftops.wasm
var pdftops []byte

//go:embed gsfonts/*
var fonts embed.FS

var fontSub fs.FS

func init() {
	sub, err := fs.Sub(fonts, "gsfonts")
	if err != nil {
		panic(err)
	}
	fontSub = sub
}

type PdfToPs struct {
	compiled wazero.CompiledModule
	wruntime wazero.Runtime
}

// DoInit inits; it is safe to call concurrently; only first will init
func New(ctx context.Context) (*PdfToPs, error) {
	r := wazero.NewRuntime(ctx)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		return nil, fmt.Errorf("cannot instantiate wasi preview 1: %w", err)
	}

	m, err := r.CompileModule(ctx, pdftops)
	if err != nil {
		return nil, fmt.Errorf("cannot compile wasm: %w", err)
	}

	if _, err := emscripten.InstantiateForModule(ctx, r, m); err != nil {
		return nil, fmt.Errorf("cannot instantiate emscriptem: %w", err)
	}

	return &PdfToPs{
		compiled: m,
		wruntime: r,
	}, nil
}

func (p *PdfToPs) ConvertPDFToPS(ctx context.Context, in []byte) ([]byte, error) {
	fs := memfs.New()

	errno := fs.WriteFile("in.pdf", in)
	if errno != 0 {
		return nil, fmt.Errorf("%s %q: %w", "cannot write initial file", "in.pdf", errno)
	}

	fsc := wazero.NewFSConfig()
	fsc = fsc.(sysfs.FSConfig).WithSysFSMount(fs, "/")

	fsc = fsc.WithFSMount(fontSub, "/usr/local/share/ghostscript/fonts/")
	fsc = fsc.WithFSMount(fontSub, "/usr/share/ghostscript/fonts/")

	stdout := new(strings.Builder)

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_start").
		WithStdout(stdout).
		WithStderr(stdout).
		WithRandSource(rand.Reader).
		WithFSConfig(fsc).
		WithName("").
		WithArgs("pdftops", "-paper", "match", "/in.pdf", "/out.ps")

	_, err := p.wruntime.InstantiateModule(ctx, p.compiled, moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot convert to PS. Output from pdftops:\n%s", stdout.String())
	}

	out, errNo := fs.ReadFile("/out.ps")
	if errNo != 0 {
		return nil, fmt.Errorf("cannot read output file:\n%d\n\noutput from pdftops:\n%s", errNo, stdout.String())
	}

	return out, nil
}
