package main

import (
	"fmt"

	"github.com/BFostek/ProvaGen/pkg/generator"
)

func main() {
	gen := generator.GeneratorFactory("")
	if err := gen.Generate("/home/breno/dev/challenges", "duplicate-integer"); err != nil {
    fmt.Println(err)
	}
}
