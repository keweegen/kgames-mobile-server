package kapp

type shutdownHandler func() error

var shutdownHandlers []shutdownHandler

func appendShutdownHandler(h shutdownHandler) {
	shutdownHandlers = append(shutdownHandlers, h)
}
