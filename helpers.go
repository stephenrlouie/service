package service

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type ServiceGroup struct {
	svcs       []Service
	wg         sync.WaitGroup
	errs       []chan error
	mergedChan chan error
	forceQuit  chan error
}

func New() *ServiceGroup {
	forceQuit := make(chan error)
	errs := make([]chan error, 1)
	return &ServiceGroup{
		errs:      errs,
		forceQuit: forceQuit,
	}
}

func (sg *ServiceGroup) Add(svc Service) {
	sg.svcs = append(sg.svcs, svc)
}

func (sg *ServiceGroup) Wait() {
	sg.wg.Wait()
}

func (sg *ServiceGroup) wrapSvc(svc Service, errs chan error) {
	defer sg.wg.Done()
	defer close(errs)
	err := svc.Start()
	if err != nil {
		errs <- err
	}
}

func (sg *ServiceGroup) Kill() {
	sg.forceQuit <- fmt.Errorf("Force Quit")
}

func (sg *ServiceGroup) Start() {
	for _, s := range sg.svcs {
		sg.wg.Add(1)
		errs := make(chan error)
		go sg.wrapSvc(s, errs)
		sg.errs = append(sg.errs, errs)
	}
	sg.mergedChan = marge(sg.errs)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	sg.wg.Add(1)
	go func() {
		sg.wg.Done()
	ctrl_loop:
		for {
			select {
			case <-signals:
				fmt.Printf("SIGINT detected. Exiting\n")
				break ctrl_loop
			case err := <-sg.mergedChan:
				if err != nil {
					fmt.Printf("Error: %v reported. Exiting\n", err)
					break ctrl_loop
				}
			case <-sg.forceQuit:
				fmt.Printf("FORCE QUIT SENDING\n")
				break ctrl_loop
			}
		}
		sg.StopAll()
	}()
}

func (sg *ServiceGroup) StopAll() {
	for _, s := range sg.svcs {
		s.Stop()
	}
}

// Marge takes multiple channels, and returns a single channel which acts as
// an aggregate. The returned channel is buffered to match the total size
// of all channel buffers.
// The returned aggregate channel will close when all originating channels
// have been closed.
func marge(errs []chan error) chan error {
	buff := 0
	for _, c := range errs {
		buff += cap(c) // cap of channel is max buffer size
	}

	// forward errors from each channel into aggregate
	wg := &sync.WaitGroup{}
	agg := make(chan error, buff)
	for _, c := range errs {
		if c != nil {
			wg.Add(1)
			go func(err <-chan error) {
				for nErr := range err {
					agg <- nErr
				}
				wg.Done()
			}(c)
		}
	}

	// close aggregate when all inputs are closed
	go func() {
		wg.Wait()
		close(agg)
	}()

	return agg
}
