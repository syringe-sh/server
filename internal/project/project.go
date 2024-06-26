package project

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/nixpig/syringe.sh/pkg/ctxkeys"
	"github.com/nixpig/syringe.sh/pkg/validation"
	"github.com/spf13/cobra"
)

func ProjectCommand() *cobra.Command {
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
		RunE:    projectAddRunE,
	}

	return projectAddCmd
}

func projectAddRunE(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	projectService, ok := cmd.Context().Value(ctxkeys.ProjectService).(ProjectService)
	if !ok {
		return fmt.Errorf("unable to get project service from context")
	}

	if err := projectService.Add(AddProjectRequest{
		Name: projectName,
	}); err != nil {
		return err
	}

	cmd.Println(fmt.Sprintf("Project '%s' added", projectName))

	return nil
}

func projectRemoveCommand() *cobra.Command {
	projectRemoveCmd := &cobra.Command{
		Use:     "remove [flags] PROJECT_NAME",
		Aliases: []string{"r"},
		Short:   "Remove a project",
		Example: "syringe project remove my_cool_project",
		Args:    cobra.MatchAll(cobra.ExactArgs(1)),
		RunE:    projectRemoveRunE,
	}

	return projectRemoveCmd
}

func projectRemoveRunE(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	projectService, ok := cmd.Context().Value(ctxkeys.ProjectService).(ProjectService)
	if !ok {
		return fmt.Errorf("unable to get project service from context")
	}

	if err := projectService.Remove(RemoveProjectRequest{
		Name: projectName,
	}); err != nil {
		return err
	}

	cmd.Println(fmt.Sprintf("Project '%s' removed", projectName))

	return nil
}

func projectRenameCommand() *cobra.Command {
	projectRenameCmd := &cobra.Command{
		Use:     "rename [flags] CURRENT_PROJECT_NAME NEW_PROJECT_NAME",
		Aliases: []string{"u"},
		Short:   "Rename a project",
		Example: "syringe project rename my_cool_project my_awesome_project",
		Args:    cobra.MatchAll(cobra.ExactArgs(2)),
		RunE:    projectRenameRunE,
	}

	return projectRenameCmd
}

func projectRenameRunE(cmd *cobra.Command, args []string) error {
	name := args[0]
	newName := args[1]

	projectService, ok := cmd.Context().Value(ctxkeys.ProjectService).(ProjectService)
	if !ok {
		return fmt.Errorf("unable to get project service from context")
	}

	if err := projectService.Rename(RenameProjectRequest{
		Name:    name,
		NewName: newName,
	}); err != nil {
		return err
	}

	cmd.Println(fmt.Sprintf("Project '%s' renamed to '%s'", name, newName))

	return nil
}

func projectListCommand() *cobra.Command {
	projectListCmd := &cobra.Command{
		Use:     "list [flags]",
		Aliases: []string{"l"},
		Short:   "List projects",
		Args:    cobra.NoArgs,
		Example: "syringe project list",
		RunE:    projectListRunE,
	}

	return projectListCmd
}

func projectListRunE(cmd *cobra.Command, args []string) error {
	projectService, ok := cmd.Context().Value(ctxkeys.ProjectService).(ProjectService)
	if !ok {
		return fmt.Errorf("unable to get project service from context")
	}

	projects, err := projectService.List()
	if err != nil {
		return err
	}

	projectNames := make([]string, len(projects.Projects))
	for i, p := range projects.Projects {
		projectNames[i] = p.Name
	}

	cmd.Print(strings.Join(projectNames, "\n"))

	return nil
}

func initProjectContext(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	db, ok := ctx.Value(ctxkeys.DB).(*sql.DB)
	if !ok {
		return fmt.Errorf("unable to get database from context")
	}

	projectService := NewProjectServiceImpl(
		NewSqliteProjectStore(db),
		validation.NewValidator(),
	)

	ctx = context.WithValue(ctx, ctxkeys.ProjectService, projectService)

	cmd.SetContext(ctx)

	return nil
}