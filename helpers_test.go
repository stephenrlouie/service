package service

import (
	"fmt"
	"testing"
	"time"
)

type testSvc struct {
	t           time.Duration
	shouldError bool
	shouldPanic bool
}

func (ts testSvc) Start() error {
	time.Sleep(ts.t)
	if ts.shouldError {
		return fmt.Errorf("Error")
	}
	if ts.shouldPanic {
		panic(fmt.Sprintf("Panic"))
	}
	return nil
}

func (ts testSvc) Stop() {}

func TestAdd(t *testing.T) {
	s := New()
	ts1 := testSvc{}
	s.Add(ts1)

	if !contains(s.svcs, ts1) {
		t.Error("Could not find ts1")
	}

	ts2 := testSvc{}
	s.Add(ts2)

	if !contains(s.svcs, ts1) {
		t.Error("Could not find ts1")
	}
	if !contains(s.svcs, ts2) {
		t.Error("Could not find ts2")
	}
}

func TestServices(t *testing.T) {
	// Working example
	tests := []struct {
		svc testSvc
	}{
		{
			// Successful
			testSvc{t: 50 * time.Millisecond, shouldError: false, shouldPanic: false},
		},
		{
			// Error
			testSvc{t: 50 * time.Millisecond, shouldError: true, shouldPanic: false},
		},
		{
			// Panic
			testSvc{t: 50 * time.Millisecond, shouldError: false, shouldPanic: true},
		},
	}

	for _, test := range tests {
		sg := New()
		sg.Add(test.svc)
		start := time.Now()
		sg.Start()
		sg.Wait()
		elapsed := time.Since(start)

		errs := sg.Status()

		if test.svc.shouldError || test.svc.shouldPanic {
			if len(errs) == 0 {
				t.Error("Expected an error or panic and received no errors")
			}
		} else {
			if len(errs) != 0 {
				t.Errorf("Expected success received %d errors", len(errs))
			}

			if elapsed < (test.svc.t) {
				t.Error("Expected to wait at least 10 seconds on service, waited less.")
			}
		}
	}
}

func TestKill(t *testing.T) {
	// Working example
	tests := []struct {
		svc      testSvc
		killTime time.Duration
	}{
		{
			// Kill clean
			svc:      testSvc{t: 100 * time.Millisecond, shouldError: false, shouldPanic: false},
			killTime: 10 * time.Millisecond,
		},
		{
			// Kill with error
			svc:      testSvc{t: 100 * time.Millisecond, shouldError: true, shouldPanic: false},
			killTime: 10 * time.Millisecond,
		},
		{
			// Kill with panic
			svc:      testSvc{t: 100 * time.Millisecond, shouldError: false, shouldPanic: true},
			killTime: 10 * time.Millisecond,
		},
		{
			// Kill after done clean
			svc:      testSvc{t: 10 * time.Millisecond, shouldError: false, shouldPanic: false},
			killTime: 100 * time.Millisecond,
		},
		{
			// Kill after done clean
			svc:      testSvc{t: 10 * time.Millisecond, shouldError: true, shouldPanic: false},
			killTime: 100 * time.Millisecond,
		},
		{
			// Kill after done clean
			svc:      testSvc{t: 50 * time.Millisecond, shouldError: false, shouldPanic: true},
			killTime: 2 * time.Second,
		},
	}

	for _, test := range tests {
		sg := New()
		sg.Add(test.svc)
		sg.Start()
		time.Sleep(test.killTime)
		sg.Kill()
		sg.Wait()

		errs := sg.Status()

		if test.svc.shouldError || test.svc.shouldPanic {
			if len(errs) == 0 {
				t.Error("Expected an error or panic and received no errors")
			}
		} else {
			if len(errs) != 0 {
				t.Errorf("Expected success received %d errors", len(errs))
			}
		}
	}
}

func contains(svcs []Service, this Service) bool {
	for _, s := range svcs {
		if this == s {
			return true
		}
	}
	return false
}
