package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/nixpig/syringe.sh/server/cmd"
	"github.com/nixpig/syringe.sh/server/internal/handlers"
	"github.com/rs/zerolog"
)

type contextKey string

const AUTHORISED = contextKey("AUTHORISED")

type SyringeSshServer struct {
	handlers handlers.SshHandlers
	log      *zerolog.Logger
}

func NewSyringeSshServer(
	handlers handlers.SshHandlers,
	log *zerolog.Logger,
) SyringeSshServer {
	return SyringeSshServer{
		handlers: handlers,
		log:      log,
	}
}

func (s SyringeSshServer) Start(host, port string) error {
	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return key.Type() == "ssh-ed25519"
		}),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					isAuthorised := sess.Context().Value(AUTHORISED)
					wish.Println(sess, fmt.Sprintf("is authorised: %v", isAuthorised))

					err := cmd.Execute(sess, s.handlers)
					if err != nil {
						wish.Error(sess, "done fucked up\n", err)
						os.Exit(1)
					}

					next(sess)
				}
			},
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					sess.Context().SetValue(
						AUTHORISED,
						s.handlers.AuthUser(sess.User(), sess.PublicKey()),
					)

					next(sess)
				}
			},
			logging.Middleware(),
		),
	)
	if err != nil {
		return err
	}

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s.log.Info().Msg("Starting SSH server")

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			s.log.Error().Err(err).Msg("Could not start server")
			done <- nil
		}
	}()

	<-done

	s.log.Info().Msg("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		s.log.Error().Err(err).Msg("Could not stop server")
	}

	return nil
}