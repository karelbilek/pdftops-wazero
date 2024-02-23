package pdftops

import (
	"context"
	"crypto/rand"
	"embed"
	"fmt"
	"io"
	"io/fs"

	"os"
	"strings"

	"github.com/karelbilek/pdftops-wazero/internal/memfs"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/experimental/sysfs"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasm/pdftops.wasm
var pdftops []byte

//go:embed gsfonts/*
var fonts embed.FS

var fontSub fs.FS

var compiled wazero.CompiledModule
var wruntime wazero.Runtime

func init() {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		panic(fmt.Errorf("cannot instantiate wasi preview 1: %w", err))
	}

	m, err := r.CompileModule(ctx, pdftops)
	if err != nil {
		panic(fmt.Errorf("cannot compile wasm: %w", err))
	}

	if _, err := emscripten.InstantiateForModule(ctx, r, m); err != nil {
		panic(fmt.Errorf("cannot instantiate emscriptem: %w", err))
	}

	sub, err := fs.Sub(fonts, "gsfonts")
	if err != nil {
		panic(err)
	}

	compiled = m
	wruntime = r
	fontSub = sub
}

func ConvertPDFToPS(ctx context.Context, in []byte) ([]byte, error) {
	fs, mfs, err := memfs.New(in)
	if err != nil {
		return nil, fmt.Errorf("cannot create a new FS: %w", err)
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
		WithArgs("pdftops", "/in.pdf", "/out.ps")

	_, err = wruntime.InstantiateModule(ctx, compiled, moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot convert to PS. Output from pdftops:\n%s", stdout.String())
	}

	outF, err := mfs.OpenFile("/out.ps", os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot open output file:\n%w\n\noutput from pdftops:\n%s", err, stdout.String())
	}
	out, err := io.ReadAll(outF)
	if err != nil {
		return nil, fmt.Errorf("cannot read output file:\n%w\n\noutput from pdftops:\n%s", err, stdout.String())
	}

	return out, nil
}
