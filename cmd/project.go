package cmd

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/nixpig/syringe.sh/server/internal/services"
	"github.com/nixpig/syringe.sh/server/internal/stores"
	"github.com/spf13/cobra"
)

const (
	projectCtxKey = contextKey("PROJECT_CTX")
)

func projectCommand() *cobra.Command {
	projectCmd := &cobra.Command{
		Use:               "project",
		Aliases:           []string{"p"},
		Short:             "Manage projects",
		PersistentPreRunE: initProjectContext,
	}

	projectCmd.AddCommand(projectAddCommand())
	projectCmd.AddCommand(projectRemoveCommand())
	projectCmd.AddCommand(projectRenameCommand())
	projectCmd.AddCommand(projectListCommand())

	return projectCmd
}

func projectAddCommand() *cobra.Command {
	projectAddCmd := &cobra.Command{
		Use:     "add [flags] PROJECT_NAME",
		Aliases: []string{"a"},
		Short:   "Add a project",
		Example: "syringe project add my_cool_project",
		Args:    cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			projectService := cmd.Context().Value(projectCtxKey).(services.ProjectService)

			if err := projectService.Add(projectName); err != nil {
				return err
			}

			return nil
		},
	}

	return projectAddCmd
}

func projectRemoveCommand() *cobra.Command {
	projectRemoveCmd := &cobra.Command{
		Use:     "remove [flags] PROJECT_NAME",
		Aliases: []string{"r"},
		Short:   "Remove a project",
		Example: "syringe project remove my_cool_project",
		Args:    cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

			projectService := cmd.Context().Value(projectCtxKey).(services.ProjectService)

			if err := projectService.Remove(projectName); err != nil {
				return err
			}

			return nil
		},
	}

	return projectRemoveCmd
}

func projectRenameCommand() *cobra.Command {
	projectRenameCmd := &cobra.Command{
		Use:     "rename [flags] CURRENT_PROJECT_NAME NEW_PROJECT_NAME",
		Aliases: []string{"r"},
		Short:   "Rename a project",
		Example: "syringe project rename my_cool_project my_awesome_project",
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			originalName := args[0]
			newName := args[1]

			projectService := cmd.Context().Value(projectCtxKey).(services.ProjectService)

			if err := projectService.Rename(originalName, newName); err != nil {
				return err
			}

			return nil
		},
	}

	return projectRenameCmd
}

func projectListCommand() *cobra.Command {
	projectListCmd := &cobra.Command{
		Use:     "list [flags]",
		Aliases: []string{"l"},
		Short:   "List projects",
		Args:    cobra.NoArgs,
		Example: "syringe project list",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectService := cmd.Context().Value(projectCtxKey).(services.ProjectService)

			projects, err := projectService.List()
			if err != nil {
				return err
			}

			if len(projects) == 0 {
				cmd.Println("No projects found!")
				cmd.Println("Try adding one with `syringe project add PROJECT_NAME`")
				return nil
			}

			for _, project := range projects {
				cmd.Println(project)
			}

			return nil
		},
	}

	return projectListCmd
}

func initProjectContext(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	db := ctx.Value(dbCtxKey).(*sql.DB)

	projectService := services.NewProjectServiceImpl(
		stores.NewSqliteProjectStore(db),
		validator.New(validator.WithRequiredStructEnabled()),
	)

	ctx = context.WithValue(ctx, projectCtxKey, projectService)

	cmd.SetContext(ctx)

	return nil
}
