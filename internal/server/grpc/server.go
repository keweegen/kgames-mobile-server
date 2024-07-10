package grpc

import (
	gamepb "github.com/keweegen/tic-toe/api/grpc/game"
	userpb "github.com/keweegen/tic-toe/api/grpc/user"
	"github.com/keweegen/tic-toe/internal/broadcaster"
	gamedomainapi "github.com/keweegen/tic-toe/internal/domain/game/api/grpc"
	userdomainapi "github.com/keweegen/tic-toe/internal/domain/user/api/grpc"
	"github.com/keweegen/tic-toe/internal/server"
	"github.com/keweegen/tic-toe/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	base        *grpc.Server
	services    *store.Service
	broadcaster broadcaster.Broadcaster
}

func NewServer(bc broadcaster.Broadcaster, services *store.Service) *Server {
	return &Server{
		broadcaster: bc,
		services:    services,
	}
}

func (s *Server) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}

	s.base = grpc.NewServer(opts...)
	s.registerHandlers()

	reflection.Register(s.base)

	return s.base.Serve(listener)
}

func (s *Server) Shutdown() error {
	s.base.GracefulStop()
	return nil
}

func (s *Server) registerHandlers() {
	gamepb.RegisterServiceServer(s.base, gamedomainapi.NewGameHandler(s.services.Game, s.services.GameStream, s.broadcaster))
	userpb.RegisterServiceServer(s.base, userdomainapi.NewUserHandler(s.services.User))
}
