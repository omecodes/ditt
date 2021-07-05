package ditt

import (
	"context"
	"io"
)

// BaseHandler defines a API calls handling pipe
type BaseHandler struct {
	Next APIHandler
}

func (b *BaseHandler) Login(ctx context.Context, user string, password string) (bool, error) {
	return b.Next.Login(ctx, user, password)
}

func (b *BaseHandler) AddUsers(ctx context.Context, reader io.Reader) error {
	return b.Next.AddUsers(ctx, reader)
}

func (b *BaseHandler) DeleteUser(ctx context.Context, userId string) error {
	return b.Next.DeleteUser(ctx, userId)
}

func (b *BaseHandler) GetUser(ctx context.Context, userId string) (UserData, error) {
	return b.Next.GetUser(ctx, userId)
}

func (b *BaseHandler) GetUserList(ctx context.Context, opts ListOptions) (*UserDataList, error) {
	return b.Next.GetUserList(ctx, opts)
}

func (b *BaseHandler) UpdateUser(ctx context.Context, userId string, userData UserData) error {
	return b.Next.UpdateUser(ctx, userId, userData)
}

// NewAPIHandler constructs an API handler pipe
func NewAPIHandler() (handler APIHandler) {

	handler = &handlerExecution{}

	handler = &handlerACL{BaseHandler: BaseHandler{
		Next: handler,
	}}

	handler = &handlerParamsValidator{
		BaseHandler{Next: handler},
	}

	return
}
