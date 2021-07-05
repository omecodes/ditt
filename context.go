package ditt

import "context"

type ctxLoggedUser struct{}

// ContextWithLoggedUser creates a new context that holds loggedUser in addition of the parent values
func ContextWithLoggedUser(parent context.Context, loggedUser string) context.Context {
	return context.WithValue(parent, ctxLoggedUser{}, loggedUser)
}

// GetLoggedUser extracts logged user name from context values
func GetLoggedUser(ctx context.Context) string {
	o := ctx.Value(ctxLoggedUser{})
	if o == nil {
		return ""
	}
	return o.(string)
}
