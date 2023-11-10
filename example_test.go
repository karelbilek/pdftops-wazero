package pdftops_test

import (
	"context"
	"fmt"
	"os"

	pdftops "github.com/karelbilek/pdftops-wazero"
)

func ExampleConvertPDFToPS() {
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

	fmt.Println("success")
	// Output: success
}
