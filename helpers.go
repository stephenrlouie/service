package service

import (
	"fmt"
	"sync"
)

type ServiceGroup struct {
	svcs       []Service
	wg         sync.WaitGroup
	errs       []chan error
	mergedChan chan error
	forceQuit  chan error
	status     []error
}

// New returns a pointer to a ServiceGroup
// This function initializes channels
func New() *ServiceGroup {
	forceQuit := make(chan error)
	errs := make([]chan error, 1)
	return &ServiceGroup{
		errs:      errs,
		forceQuit: forceQuit,
	}
}

// Add will take a service and Add it to the ServiceGroup
func (sg *ServiceGroup) Add(svc Service) {
	sg.svcs = append(sg.svcs, svc)
}

// Wait wraps sync.WaitGroup.Wait() and will ensure all children
// routines in the ServiceGroup conclude before the parent process moves on
func (sg *ServiceGroup) Wait() {
	sg.wg.Wait()
}

// Kill is a way for the parent to force all children routines in
// the ServiceGroup to stop
func (sg *ServiceGroup) Kill() {
	sg.forceQuit <- fmt.Errorf("Force Quit")
}

// Start will begin every child routine in the ServiceGroup
// and listen on the error channels for the children routines.
// If an error is received it will close all other routines in the
// ServiceGroup
func (sg *ServiceGroup) Start() {
	for _, s := range sg.svcs {
		sg.wg.Add(1)
		errs := make(chan error)
		go sg.wrapSvc(s, errs)
		sg.errs = append(sg.errs, errs)
	}
	sg.mergedChan = merge(sg.errs)

	sg.wg.Add(1)
	go func() {
		defer sg.wg.Done()
	ctrl_loop:
		for {
			select {
			case err := <-sg.mergedChan:
				if err != nil {
					sg.status = append(sg.status, err)
					break ctrl_loop
				}
			case <-sg.forceQuit:
				break ctrl_loop
			}
		}
		sg.stopAll()

		// Receive any final shutdown errors
		for _, errCh := range sg.errs {
			if errCh != nil {
				for err := range errCh {
					if err != nil {
						sg.status = append(sg.status, err)
					}
				}
			}
		}
	}()
}

func (sg *ServiceGroup) Status() []error {
	return sg.status
}

func (sg *ServiceGroup) wrapSvc(svc Service, errs chan error) {
	defer sg.wg.Done()
	defer close(errs)
	err := svc.Start()
	if err != nil {
		errs <- err
	}
}

func (sg *ServiceGroup) stopAll() {
	for _, s := range sg.svcs {
		s.Stop()
	}
}

// Marge takes multiple channels, and returns a single channel which acts as
// an aggregate. The returned channel is buffered to match the total size
// of all channel buffers.
// The returned aggregate channel will close when all originating channels
// have been closed.
func merge(errs []chan error) chan error {
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
