package main

import (
	"github.com/BFostek/ProvaGen/pkg/generator"
)

func main() {
  gen := generator.GeneratorFactory("")
  gen.Generate("", "duplicate-integer")
  
}
