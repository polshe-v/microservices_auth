package model

// EndpointPermissions type is the structure for endpoint permissions by roles.
type EndpointPermissions struct {
	Endpoint string
	Roles    []string
}
