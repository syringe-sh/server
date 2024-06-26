package environment

import (
	"github.com/nixpig/syringe.sh/pkg/serrors"
	"github.com/nixpig/syringe.sh/pkg/validation"
)

type AddEnvironmentRequest struct {
	Name    string `name:"environment name" validate:"required,min=1,max=256"`
	Project string `name:"project name" validate:"required,min=1,max=256"`
}

type RemoveEnvironmentRequest struct {
	Name    string `name:"environment name" validate:"required,min=1,max=256"`
	Project string `name:"project name" validate:"required,min=1,max=256"`
}

type RenameEnvironmentRequest struct {
	Name    string `name:"environment name" validate:"required,min=1,max=256"`
	NewName string `name:"new environment name" validate:"required,min=1,max=256"`
	Project string `name:"project name" validate:"required,min=1,max=256"`
}

type ListEnvironmentRequest struct {
	Project string `name:"project name" validate:"required,min=1,max=256"`
}

type EnvironmentResponse struct {
	ID   int
	Name string
}

type ListEnvironmentsResponse struct {
	Project      string
	Environments []EnvironmentResponse
}

type EnvironmentService interface {
	Add(environment AddEnvironmentRequest) error
	Remove(environment RemoveEnvironmentRequest) error
	Rename(environment RenameEnvironmentRequest) error
	List(project ListEnvironmentRequest) (*ListEnvironmentsResponse, error)
}

func NewEnvironmentServiceImpl(
	store EnvironmentStore,
	validate validation.Validator,
) EnvironmentService {
	return EnvironmentServiceImpl{
		store:    store,
		validate: validate,
	}
}

type EnvironmentServiceImpl struct {
	store    EnvironmentStore
	validate validation.Validator
}

func (e EnvironmentServiceImpl) Add(
	environment AddEnvironmentRequest,
) error {
	if err := e.validate.Struct(environment); err != nil {
		return serrors.ValidationError(err)
	}

	if err := e.store.Add(
		environment.Name,
		environment.Project,
	); err != nil {
		return err
	}

	return nil
}

func (e EnvironmentServiceImpl) Remove(
	environment RemoveEnvironmentRequest,
) error {
	if err := e.validate.Struct(environment); err != nil {
		return serrors.ValidationError(err)
	}

	if err := e.store.Remove(
		environment.Name,
		environment.Project,
	); err != nil {
		return err
	}

	return nil
}

func (e EnvironmentServiceImpl) Rename(
	environment RenameEnvironmentRequest,
) error {
	if err := e.validate.Struct(environment); err != nil {
		return serrors.ValidationError(err)
	}

	if err := e.store.Rename(
		environment.Name,
		environment.NewName,
		environment.Project,
	); err != nil {
		return err
	}

	return nil
}

func (e EnvironmentServiceImpl) List(
	request ListEnvironmentRequest,
) (*ListEnvironmentsResponse, error) {
	if err := e.validate.Struct(request); err != nil {
		return nil, serrors.ValidationError(err)
	}

	environments, err := e.store.List(request.Project)
	if err != nil {
		return nil, err
	}

	var environmentsResponseList []EnvironmentResponse

	for _, ev := range *environments {
		environmentsResponseList = append(environmentsResponseList, EnvironmentResponse{
			ID:   ev.ID,
			Name: ev.Name,
		},
		)
	}

	return &ListEnvironmentsResponse{
		Project:      request.Project,
		Environments: environmentsResponseList,
	}, nil
}
