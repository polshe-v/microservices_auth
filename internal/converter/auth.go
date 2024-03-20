package converter

import (
	"github.com/polshe-v/microservices_auth/internal/model"
	desc "github.com/polshe-v/microservices_auth/pkg/auth_v1"
)

// ToUserLoginFromDesc converts structure of API layer to service layer model.
func ToUserLoginFromDesc(creds *desc.Creds) *model.UserCreds {
	return &model.UserCreds{
		Username: creds.Username,
		Password: creds.Password,
	}
}
