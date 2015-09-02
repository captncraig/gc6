package main

import (
	"fmt"
	"github.com/golangchallenge/gc6/generators"
)

func main() {
	m := generators.DepthFirst()
	fmt.Println(m)
}
