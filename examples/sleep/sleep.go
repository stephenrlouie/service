package sleep

import (
	"fmt"
	"time"
)

type Sleep struct {
	Msg     string
	Pass    bool
	Seconds int
	Quit    bool
}

func (s *Sleep) Start(errs chan error) {
	defer close(errs)
	s.Quit = false
	ticker := time.NewTicker(time.Second * 1)
	count := 0
ctrl_loop:
	for {
		select {
		case <-ticker.C:
			if s.Quit || count == s.Seconds {
				fmt.Printf("%s exit loop entered\n", s.Msg)
				break ctrl_loop
			}

			count++
			fmt.Printf("%s!\n", s.Msg)
		}
	}

	if !s.Pass {
		errs <- fmt.Errorf("TOO MUCH SLEEP")
	}
}

func (s *Sleep) Stop() {
	fmt.Printf("%s stop\n", s.Msg)
	s.Quit = true
}
