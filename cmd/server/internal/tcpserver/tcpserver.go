package tcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"golang.org/x/sync/errgroup"
)

type Handler interface {
	HandleConnection(conn net.Conn) error
}

type QuoteServer struct {
	address string
	handler Handler
}

func New(
	address string,
	handler Handler,
) *QuoteServer {
	return &QuoteServer{
		address: address,
		handler: handler,
	}
}

func (s *QuoteServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("create tcp listener: %w", err)
	}
	defer func() {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			slog.Error("close tcp listener", "error", err)
		}
	}()

	slog.Debug("server started", "address", s.address)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
					return nil
				}

				return fmt.Errorf("accept connection: %w", err)
			}

			g.Go(func() error {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("connection handler panic", "error", fmt.Errorf("%v", r))
					}
				}()

				defer func() {
					if err := conn.Close(); err != nil {
						slog.Error("close connection", "error", err)
					}
				}()

				if err := s.handler.HandleConnection(conn); err != nil {
					slog.Error("handle connection", "error", err)
				}

				return nil
			})
		}
	})

	<-ctx.Done()

	if err := listener.Close(); err != nil {
		slog.Error("close tcp listener", "error", err)
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("wait for connections: %w", err)
	}

	slog.Debug("server shut down")

	return nil
}
