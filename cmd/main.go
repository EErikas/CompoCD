package main

import (
	"fmt"

	"github.com/EErikas/CompoCD/internal/stack"
)

func main() {
	var s stack.Node
	if err := s.ParseStackConfig("stack.yaml"); err == nil {
		fmt.Printf("%+v\n", s)
		s.StartStack()
	}
}
