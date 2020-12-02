package mocks

type AuthenticationServiceMockInterface interface {
	SetValidatePassword(f func(plainPassword []byte, hashedPassword string) (bool, error))
}

type AuthenticationServiceMock struct {
	validatePassword func(plainPassword []byte, hashedPassword string) (bool, error)
}

func (m *AuthenticationServiceMock) SetValidatePassword(f func(plainPassword []byte, hashedPassword string) (bool, error)) {
	m.validatePassword = f
}

func (m *AuthenticationServiceMock) ValidatePassword(plainPassword []byte, hashedPassword string) (bool, error) {
	return m.validatePassword(plainPassword, hashedPassword)
}
