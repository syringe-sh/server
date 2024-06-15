package services

import (
	"github.com/go-playground/validator/v10"
	"github.com/nixpig/syringe.sh/server/internal/stores"
)

type AddProjectRequest struct {
	Name string `validate:"required,min=1,max=256"`
}

type RemoveProjectRequest struct {
	Name string `validate:"required,min=1,max=256"`
}

type RenameProjectRequest struct {
	Name    string `validate:"required,min=1,max=256"`
	NewName string `validate:"required,min=1,max=256"`
}

type ProjectService interface {
	Add(project AddProjectRequest) error
	Remove(project RemoveProjectRequest) error
	Rename(project RenameProjectRequest) error
	List() ([]string, error)
}

func NewProjectServiceImpl(
	store stores.ProjectStore,
	validate *validator.Validate,
) ProjectService {
	return ProjectServiceImpl{
		store:    store,
		validate: validate,
	}
}

type ProjectServiceImpl struct {
	store    stores.ProjectStore
	validate *validator.Validate
}

func (p ProjectServiceImpl) Add(project AddProjectRequest) error {
	if err := p.validate.Struct(project); err != nil {
		return err
	}

	if err := p.store.Add(project.Name); err != nil {
		return err
	}

	return nil
}

func (p ProjectServiceImpl) Remove(project RemoveProjectRequest) error {
	if err := p.validate.Struct(project); err != nil {
		return err
	}

	if err := p.store.Remove(project.Name); err != nil {
		return err
	}

	return nil
}

func (p ProjectServiceImpl) Rename(project RenameProjectRequest) error {
	if err := p.validate.Struct(project); err != nil {
		return err
	}

	if err := p.store.Rename(
		project.Name,
		project.NewName,
	); err != nil {
		return err
	}

	return nil
}

func (p ProjectServiceImpl) List() ([]string, error) {
	projects, err := p.store.List()
	if err != nil {
		return nil, err
	}

	return projects, nil
}