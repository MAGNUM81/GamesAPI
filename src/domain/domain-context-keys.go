package domain

type contextKey string

var (
	contextKeyRbacUserId = contextKey("rbac.userId")
)

func (c contextKey) String() string {
	return "domain context key " + string(c)
}

func ContextKey(s string) string {
	return contextKey(s).String()
}

func RbacUserId() string {
	return contextKeyRbacUserId.String()
}
