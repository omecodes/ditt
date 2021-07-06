package ditt

import (
	"context"
)

type handlerACL struct {
	BaseHandler
}

func (h *handlerACL) assertIsAdmin(ctx context.Context) error {
	loggedUser := GetLoggedUser(ctx)
	if loggedUser != "admin" {
		return Forbidden
	}
	return nil
}

func (h *handlerACL) assertIsAuthenticated(ctx context.Context) error {
	if GetLoggedUser(ctx) == "" {
		return Forbidden
	}
	return nil
}

func (h *handlerACL) assertHasAccess(ctx context.Context, userId string) error {
	loggedUser := GetLoggedUser(ctx)
	if loggedUser == "" {
		return Forbidden
	}

	if loggedUser != userId && loggedUser != "admin" {
		return NotAuthorized
	}

	return nil
}

func (h *handlerACL) Login(ctx context.Context, login string, password string) (bool, error) {
	if err := h.assertIsAuthenticated(ctx); err == nil {
		return false, NotAuthorized
	}
	return h.BaseHandler.Login(ctx, login, password)
}

/* func (h *handlerACL) AddUsers(ctx context.Context, reader io.Reader) error {
	err := h.assertIsAdmin(ctx)
	if err != nil {
		return err
	}
	return h.BaseHandler.AddUsers(ctx, reader)
} */

func (h *handlerACL) DeleteUser(ctx context.Context, userId string) error {
	err := h.assertHasAccess(ctx, userId)
	if err != nil {
		return err
	}

	return h.BaseHandler.DeleteUser(ctx, userId)
}

func (h *handlerACL) GetUser(ctx context.Context, userId string) (UserData, error) {
	err := h.assertHasAccess(ctx, userId)
	if err != nil {
		return "", err
	}

	return h.BaseHandler.GetUser(ctx, userId)
}

func (h *handlerACL) GetUserList(ctx context.Context, opts ListOptions) (*UserDataList, error) {
	err := h.assertIsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	return h.BaseHandler.GetUserList(ctx, opts)
}

func (h *handlerACL) UpdateUser(ctx context.Context, userId string, userData UserData) error {
	err := h.assertHasAccess(ctx, userId)
	if err != nil {
		return err
	}

	return h.BaseHandler.UpdateUser(ctx, userId, userData)
}
