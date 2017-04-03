package service

// A Service is a long-running cancellable routine.
// Services are responsible for closing their errs channel
// and should call Done on the provided WaitGroup on shutdown
type Service interface {
	Start() error
	Stop()
}
