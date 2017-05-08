package sleep

import (
	"fmt"
	"time"
)

type Sleep struct {
	Id      string
	Pass    bool
	Panic   bool
	Seconds int
	Quit    bool
}

func (s *Sleep) Start() error {
	defer fmt.Printf("%s is closed\n", s.Id)
	s.Quit = false
	ticker := time.NewTicker(time.Second * 1)
	count := 0
ctrl_loop:
	for {
		select {
		case <-ticker.C:
			if s.Quit || count == s.Seconds {
				break ctrl_loop
			}

			count++
			fmt.Printf("%s!\n", s.Id)
		}
	}

	if !s.Pass {
		return fmt.Errorf("%s fail", s.Id)
	}
	if s.Panic {
		panic(fmt.Sprintf("%s panic", s.Id))
	}
	return nil
}

func (s *Sleep) Stop() {
	fmt.Printf("%s stop\n", s.Id)
	s.Quit = true
}
