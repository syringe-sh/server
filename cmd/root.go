package cmd

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/ssh"
	"github.com/nixpig/syringe.sh/server/internal/database"
	"github.com/nixpig/syringe.sh/server/internal/services"
	"github.com/nixpig/syringe.sh/server/pkg/turso"
	"github.com/spf13/cobra"
	gossh "golang.org/x/crypto/ssh"
)

type contextKey string

const (
	dbCtxKey   = contextKey("DB_CTX")
	sessCtxKey = contextKey("SESS_CTX")
)

func Execute(
	sess ssh.Session,
	appService services.AppService,
) error {
	rootCmd := &cobra.Command{
		Use:   "syringe",
		Short: "Distributed environment variable management over SSH.",
		Long:  "Distributed environment variable management over SSH.",
	}

	rootCmd.AddCommand(userCommand(sess, appService))

	rootCmd.AddCommand(projectCommand())
	rootCmd.AddCommand(environmentCommand())
	rootCmd.AddCommand(secretCommand())

	rootCmd.SetArgs(sess.Command())
	rootCmd.SetIn(sess)
	rootCmd.SetOut(sess)
	rootCmd.SetErr(sess.Stderr())
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	ctx := context.Background()

	ctx = context.WithValue(ctx, sessCtxKey, sess)

	db, err := NewUserDB(sess.PublicKey())
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, dbCtxKey, db)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	return nil
}

// TODO: really don't like this!!
func NewUserDB(publicKey ssh.PublicKey) (*sql.DB, error) {
	api := turso.New(
		os.Getenv("DATABASE_ORG"),
		os.Getenv("API_TOKEN"),
		http.Client{},
	)

	marshalledKey := gossh.MarshalAuthorizedKey(publicKey)

	hashedKey := fmt.Sprintf("%x", sha1.Sum(marshalledKey))
	expiration := "30s"

	token, err := api.CreateToken(hashedKey, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to create token:\n%s", err)
	}

	fmt.Println("creating new user-specific db connection")
	db, err := database.Connection(
		"libsql://"+hashedKey+"-"+os.Getenv("DATABASE_ORG")+".turso.io",
		string(token.Jwt),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection:\n%s", err)
	}

	return db, nil
}
