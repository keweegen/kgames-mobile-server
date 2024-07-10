package server

type Server interface {
	Run(addr string) error
	Shutdown() error
}
