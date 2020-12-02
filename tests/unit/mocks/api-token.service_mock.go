package mocks

type TokenServiceMockInterface interface {
	SetValidateToken(func(string) (bool, error))
}
type TokenServiceMock struct {
	validateToken func(string) (bool, error)
}

func (t *TokenServiceMock) SetValidateToken(f func(string) (bool, error)) {
	t.validateToken = f
}

func (t *TokenServiceMock) GetApiToken() (token string, err error) {
	return "1234", nil
}

func (t *TokenServiceMock) ValidateToken(s string) (token bool, err error) {
	return t.validateToken(s)
}
