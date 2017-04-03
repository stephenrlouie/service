package hello

import "fmt"

type Hello struct{}

func (h *Hello) Start(errs chan error) {
	defer close(errs)
	fmt.Printf("Hello world\n")
}

func (h *Hello) Stop() {
	fmt.Printf("Hello stop\n")
}
