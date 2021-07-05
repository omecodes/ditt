package ditt

import (
	"context"
	"io"
)

type handlerParamsValidator struct {
	BaseHandler
}

func (h handlerParamsValidator) Login(ctx context.Context, login string, password string) (bool, error) {
	if login == "" {
		return false, BadInput
	}

	if password == "" {
		return false, BadInput
	}

	return h.BaseHandler.Login(ctx, login, password)
}

func (h handlerParamsValidator) AddUsers(ctx context.Context, reader io.Reader) error {
	if reader == nil {
		return BadInput
	}
	return h.BaseHandler.AddUsers(ctx, reader)
}

func (h handlerParamsValidator) DeleteUser(ctx context.Context, userId string) error {
	if userId == "" {
		return BadInput
	}
	return h.BaseHandler.DeleteUser(ctx, userId)
}

func (h handlerParamsValidator) GetUser(ctx context.Context, userId string) (UserData, error) {
	if userId == "" {
		return "", BadInput
	}
	return h.BaseHandler.GetUser(ctx, userId)
}

func (h handlerParamsValidator) GetUserList(ctx context.Context, opts ListOptions) (*UserDataList, error) {
	if opts.Offset < 0 {
		return nil, BadInput
	}

	if opts.Count < 0 {
		return nil, BadInput
	}

	if opts.Count == 0 || opts.Count > DefaultUserListCount {
		opts.Count = DefaultUserListCount
	}

	return h.BaseHandler.GetUserList(ctx, opts)
}

func (h handlerParamsValidator) UpdateUser(ctx context.Context, userId string, userData UserData) error {
	if userId == "" {
		return BadInput
	}

	if userData == "" {
		return BadInput
	}

	return h.BaseHandler.UpdateUser(ctx, userId, userData)
}
