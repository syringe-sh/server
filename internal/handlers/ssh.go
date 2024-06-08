package handlers

import (
	"os"

	"github.com/charmbracelet/ssh"
	"github.com/nixpig/syringe.sh/server/internal/services"
	"github.com/rs/zerolog"
)

type SshHandlers struct {
	appService services.AppService
	log        *zerolog.Logger
}

func NewSshHandlers(appService services.AppService, log *zerolog.Logger) SshHandlers {
	return SshHandlers{
		appService: appService,
		log:        log,
	}
}

func (h *SshHandlers) RegisterUser(username string, publicKey ssh.PublicKey) {
	registeredUser, err := h.appService.RegisterUser(services.RegisterUserRequest{
		Username:  username,
		PublicKey: publicKey,
		Email:     "notusedyet@example.org",
	})
	if err != nil {
		h.log.Error().Err(err).Msg("failed to register with error")
		os.Exit(1)
	}

	h.log.Info().Str("username", registeredUser.Username).Msg("registered user")
}

func (h *SshHandlers) AuthUser(
	username string,
	publicKey ssh.PublicKey,
) bool {
	userAuth, err := h.appService.AuthenticateUser(services.UserAuthRequest{
		Username:  username,
		PublicKey: publicKey,
	})
	if err != nil {
		h.log.Error().Err(err).Msg("failed auth with error")
		os.Exit(1)
	}

	return userAuth.Auth
}
