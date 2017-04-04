package hello

import "fmt"

type Hello struct {
	Id string
}

func (h *Hello) Start() error {
	fmt.Printf("%s says: 'Hello world'\n", h.Id)
	return nil
}

func (h *Hello) Stop() {
	fmt.Printf("%s stop\n", h.Id)
}
