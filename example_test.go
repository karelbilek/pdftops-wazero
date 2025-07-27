package pdftopswazero_test

import (
	"context"
	"fmt"
	"os"

	pdftops "github.com/karelbilek/pdftopswazero"
)

func Example_convertPDFToPS() {
	g, err := pdftops.New(context.Background())
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

	fmt.Println("success")
	// Output: success
}
