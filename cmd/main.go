package main

import (
	"fmt"
	"time"

	"github.com/EErikas/CompoCD/internal/stack"
)

func main() {
	var s stack.Node
	if err := s.ParseStackConfig("stack.yaml"); err == nil {
		loop(s)
	}
}

func loop(s stack.Node) {
	for {
		if err := s.UpdateRepo(); err != nil {
			fmt.Printf("Error %e\n", err)
		}
		// if err := s.StartStack(); err != nil {
		// 	fmt.Printf("Error %e\n", err)
		// }
		s.StartStack()
		time.Sleep(time.Second * 5)
	}
}
