package converter

import (
	model "github.com/polshe-v/microservices_auth/internal/model"
	modelRepo "github.com/polshe-v/microservices_auth/internal/repository/access/model"
)

// ToEndpointPermissionsFromRepo converts repository layer model to structure of service layer.
func ToEndpointPermissionsFromRepo(endpointPermissions []*modelRepo.EndpointPermissions) []*model.EndpointPermissions {
	var res []*model.EndpointPermissions
	for _, e := range endpointPermissions {
		res = append(res, &model.EndpointPermissions{
			Endpoint: e.Endpoint,
			Roles:    e.Roles,
		})
	}
	return res
}
