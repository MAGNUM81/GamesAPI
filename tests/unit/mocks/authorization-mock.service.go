package mocks

import (
	"GamesAPI/src/domain"
	"context"
	"net/url"
)

type AuthorizationServiceMockInterface interface {
	SetAuthorize(f func(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error)
}

type AuthorizationServiceMock struct {
	authorize func(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error
}

func (s *AuthorizationServiceMock) GetRbac() domain.RBAC {
	return nil
}

func (s *AuthorizationServiceMock) SetAuthorize(f func(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error) {
	s.authorize = f
}

func (s *AuthorizationServiceMock) Authorize(ctx context.Context, url *url.URL, role string, resource string, endpoint string) error {
	return s.authorize(ctx, url, role, resource, endpoint)
}
