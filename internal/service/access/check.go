package access

import (
	"context"
	"errors"
	"log"
	"slices"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/polshe-v/microservices_auth/internal/converter"
	"github.com/polshe-v/microservices_auth/internal/utils"
)

const (
	authMetadataHeader = "authorization"
	authPrefix         = "Bearer "
)

var accessibleRoles map[string][]string

func (s *serv) Check(ctx context.Context, endpoint string) error {
	var err error
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md[authMetadataHeader]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	// Get secret key from storage for access token HMAC
	accessTokenSecretKey, err := s.keyRepository.GetKey(ctx, accessTokenSecretKeyName)
	if err != nil {
		log.Print(err)
		return errors.New("failed to generate token")
	}

	if accessibleRoles == nil {
		endpointPermissions, errRepo := s.accessRepository.GetRoleEndpoints(ctx)
		if errRepo != nil {
			log.Print(errRepo)
			return errors.New("failed to read access policy")
		}
		accessibleRoles = converter.ToEndpointPermissionsMap(endpointPermissions)
	}

	// Read slice of roles allowed to use the endpoint
	roles, ok := accessibleRoles[endpoint]
	if !ok {
		log.Print(err)
		return errors.New("failed to find endpoint")
	}

	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
	if err != nil {
		return errors.New("access token is invalid")
	}

	// If role is not in the slice of roles allowed to use the endpoint, then deny access
	if !slices.Contains(roles, claims.Role) {
		return errors.New("access denied")
	}

	return nil
}
