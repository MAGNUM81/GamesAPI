package services

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"context"
	"net/url"
	"os"
)

var (
	AuthorizationService AuthorizationServiceInterface = NewAuthorizationService(os.Getenv("RBAC_FILEPATH"))
)

type AuthorizationServiceInterface interface {
	Authorize(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error
}

type authorizationService struct {
	rbac domain.RBAC
}

//Constructor - must be instantiated with a role-based access YAML file
func NewAuthorizationService(path string) *authorizationService {
	ret := &authorizationService{}
	rbac, err := domain.RbacFromFile(path)
	if err != nil {
		panic(err)
	}
	ret.rbac = rbac
	return ret
}

func (a authorizationService) Authorize(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error {
	permission, exists := a.rbac[role][resource][endpoint]

	if !exists {
		return errorUtils.ErrRoleUnknown
	}

	if !permission.Allow {
		return errorUtils.ErrForbidden
	}

	err := permission.Ensure.QueryComplies(ctx, url)
	if err != nil {
		return err
	}

	err = permission.Ensure.QueryComplies(ctx, url)
	if err != nil {
		return err
	}

	return nil
}







