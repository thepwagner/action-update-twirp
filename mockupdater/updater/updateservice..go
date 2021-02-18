package updater

import (
	"context"

	v1 "github.com/thepwagner/action-update-twirp/proto/actionupdate/v1"
)

type UpdateService struct {
}

func NewUpdateService() *UpdateService {
	return &UpdateService{}
}

var _ v1.UpdateService = (*UpdateService)(nil)

func (u *UpdateService) ListDependencies(context.Context, *v1.ListDependenciesRequest) (*v1.ListDependenciesResponse, error) {
	return &v1.ListDependenciesResponse{
		Dependencies: []*v1.Dependency{
			{Path: "foo", Version: "v1.0.0"},
			{Path: "bar", Version: "v1.0.0"},
		},
	}, nil
}
