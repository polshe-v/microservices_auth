package converter

import (
	"github.com/polshe-v/microservices_auth/internal/model"
)

// ToEndpointPermissionsMap converts slice of service layer structures to map.
func ToEndpointPermissionsMap(endpointPermissions []*model.EndpointPermissions) map[string][]string {
	res := make(map[string][]string)
	for _, e := range endpointPermissions {
		res[e.Endpoint] = e.Roles
	}
	return res
}
