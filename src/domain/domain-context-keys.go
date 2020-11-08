package domain

type contextKey string

var (
	contextKeyRbacEmail = contextKey("rbac.email")
	contextKeyRbacAccess = contextKey("rbac.access")
)

func (c contextKey) String() string {
	return "domain context key " + string(c)
}

func ContextKey(s string) string {
	return contextKey(s).String()
}

func RbacEmail() string {
	return contextKeyRbacEmail.String()
}

func RbacAccess() string {
	return contextKeyRbacAccess.String()
}
