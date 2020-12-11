package domain

type contextKey string

var (
	contextKeyRbacUserId = contextKey("userId")
)

func (c contextKey) String() string {
	return string(c)
}

func ContextKey(s string) string {
	return contextKey(s).String()
}

func RbacUserId() string {
	return contextKeyRbacUserId.String()
}
