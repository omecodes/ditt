package ditt

import (
	"bytes"
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseHandler_AddUsers1(t *testing.T) {
	Convey("Calling AddUsers with a nil stream must fail", t, func() {
		handler := NewAPIHandler()
		err := handler.AddUsers(context.Background(), nil)
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_AddUsers2(t *testing.T) {
	Convey("Calling AddUsers with an unauthenticated context must fail", t, func() {
		handler := NewAPIHandler()
		err := handler.AddUsers(context.Background(), bytes.NewBufferString(""))
		So(err, ShouldEqual, Forbidden)
	})
}

func TestBaseHandler_AddUsers3(t *testing.T) {
	Convey("Calling AddUsers with a non admin context must fail", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "user1")
		err := handler.AddUsers(authenticatedContext, bytes.NewBufferString(""))
		So(err, ShouldEqual, Forbidden)
	})
}

func TestBaseHandler_AddUsers4(t *testing.T) {
	Convey("Calling AddUsers with an admin context and with malformed JSON must fail", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "admin")

		userDataStream := `
			[
				{"id": "loki",  "data": "lorem ipsum"},
				{"id": "hulk", "data": lorem ipsum"
			]
		`
		err := handler.AddUsers(authenticatedContext, bytes.NewBufferString(userDataStream))
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_AddUsers5(t *testing.T) {
	Convey("Calling AddUsers with an admin context and with a well formed JSON should succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "admin")
		userDataStream := `
			[
				{"id": "loki", "password": "loki-pass", "data": "lorem ipsum"},
				{"id": "hulk", "password": "hulk-pass", "data": "lorem ipsum"}
			]
		`
		err := handler.AddUsers(authenticatedContext, bytes.NewBufferString(userDataStream))
		So(err, ShouldBeNil)
	})
}

func TestBaseHandler_Login1(t *testing.T) {
	Convey("Calling Login with an empty login or an empty password must fail", t, func() {
		handler := NewAPIHandler()
		ok, err := handler.Login(context.Background(), "", "password")
		So(ok, ShouldBeFalse)
		So(err, ShouldEqual, BadInput)

		ok, err = handler.Login(context.Background(), "login", "")
		So(ok, ShouldBeFalse)
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_Login2(t *testing.T) {
	Convey("Calling Login with an authenticated context must fail", t, func() {
		handler := NewAPIHandler()

		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		ok, err := handler.Login(authenticatedContext, "loki", "password")
		So(ok, ShouldBeFalse)
		So(err, ShouldEqual, NotAuthorized)
	})
}

func TestBaseHandler_Login3(t *testing.T) {
	Convey("Calling Login with an unauthenticated context and with a userId that does not exists must fail", t, func() {
		handler := NewAPIHandler()
		ok, err := handler.Login(context.Background(), "bamba", "loki-pass")
		So(ok, ShouldBeFalse)
		So(err, ShouldBeNil)
	})
}

func TestBaseHandler_Login4(t *testing.T) {
	Convey("Calling Login with an unauthenticated context and with correct credentials should succeed", t, func() {
		handler := NewAPIHandler()
		ok, err := handler.Login(context.Background(), "loki", "loki-pass")
		So(ok, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

func TestBaseHandler_GetUser1(t *testing.T) {
	Convey("Calling GetUser with an empty userId must fail", t, func() {
		handler := NewAPIHandler()
		_, err := handler.GetUser(context.Background(), "")
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_GetUser2(t *testing.T) {
	Convey("Calling GetUser with an unauthenticated context must fail", t, func() {
		handler := NewAPIHandler()
		_, err := handler.GetUser(context.Background(), "user1")
		So(err, ShouldEqual, Forbidden)
	})
}

func TestBaseHandler_GetUser3(t *testing.T) {
	Convey("Calling GetUser with an authenticated context on another user data must fail", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		_, err := handler.GetUser(authenticatedContext, "hulk")
		So(err, ShouldEqual, NotAuthorized)
	})
}

func TestBaseHandler_GetUser4(t *testing.T) {
	Convey("Calling GetUser with an authenticated context on owned data must succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "admin")

		userDataStream := `
			[
				{"id": "loki",  "data": "lorem ipsum"},
				{"id": "hulk", "data": lorem ipsum"
			]
		`
		err := handler.AddUsers(authenticatedContext, bytes.NewBufferString(userDataStream))
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_GetUser5(t *testing.T) {
	Convey("Calling AddUsers with an admin context and with a well formed JSON should succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		userData, err := handler.GetUser(authenticatedContext, "loki")
		So(err, ShouldBeNil)
		So(userData, ShouldNotEqual, "")
		So(userData.Id(), ShouldEqual, "loki")
	})
}

func TestBaseHandler_GetUserList1(t *testing.T) {
	Convey("Calling GetUserList with negative bounds must fail", t, func() {
		handler := NewAPIHandler()
		_, err := handler.GetUserList(context.Background(), ListOptions{
			Offset: -1,
		})
		So(err, ShouldEqual, BadInput)

		_, err = handler.GetUserList(context.Background(), ListOptions{
			Count: -1,
		})
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_GetUserList2(t *testing.T) {
	Convey("Calling GetUserList with an unauthenticated context must fail", t, func() {
		handler := NewAPIHandler()
		_, err := handler.GetUserList(context.Background(), ListOptions{})
		So(err, ShouldEqual, Forbidden)
	})
}

func TestBaseHandler_GetUserList3(t *testing.T) {
	Convey("Calling GetUserList with an authenticated context should succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		list, err := handler.GetUserList(authenticatedContext, ListOptions{})
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)

		for _, userData := range list.UserDataList {
			So(userData.Id(), ShouldEqual, "loki")
		}
	})
}

func TestBaseHandler_GetUserList4(t *testing.T) {
	Convey("Calling GetUserList with an admin context should succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "admin")
		list, err := handler.GetUserList(authenticatedContext, ListOptions{})
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
	})
}

func TestBaseHandler_UpdateUser1(t *testing.T) {
	Convey("Calling UpdateUser with an empty userId or userData must fail", t, func() {
		handler := NewAPIHandler()
		err := handler.UpdateUser(context.Background(), "", "whatever")
		So(err, ShouldEqual, BadInput)

		err = handler.UpdateUser(context.Background(), "user1", "")
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_UpdateUser2(t *testing.T) {
	Convey("Calling UpdateUser with an unauthenticated context must fail", t, func() {
		handler := NewAPIHandler()
		err := handler.UpdateUser(context.Background(), "loki", "whatever")
		So(err, ShouldEqual, Forbidden)
	})
}

func TestBaseHandler_UpdateUser3(t *testing.T) {
	Convey("Calling UpdateUser with an authenticated context on another user data must fail", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		err := handler.UpdateUser(authenticatedContext, "hulk", "whatever")
		So(err, ShouldEqual, NotAuthorized)
	})
}

func TestBaseHandler_UpdateUser4(t *testing.T) {
	Convey("Calling UpdateUser with an authenticated context on owned data must succeed", t, func() {
		handler := NewAPIHandler()
		userData := UserData(`{"id": "loki", "data": I am a god you dummy creatures"}`)
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		err := handler.UpdateUser(authenticatedContext, "loki", userData)
		So(err, ShouldBeNil)
	})
}

func TestBaseHandler_UpdateUser5(t *testing.T) {
	Convey("Calling UpdateUser with an admin context must succeed", t, func() {
		handler := NewAPIHandler()
		userData := UserData(`{"id": "hulk", "data": I don't have time to think, all i want to destroy you"}`)
		authenticatedContext := ContextWithLoggedUser(context.Background(), "admin")
		err := handler.UpdateUser(authenticatedContext, "hulk", userData)
		So(err, ShouldBeNil)
	})
}

func TestBaseHandler_DeleteUser1(t *testing.T) {
	Convey("Calling DeleteUser with an empty userId must fail", t, func() {
		handler := NewAPIHandler()
		err := handler.DeleteUser(context.Background(), "")
		So(err, ShouldEqual, BadInput)
	})
}

func TestBaseHandler_DeleteUser2(t *testing.T) {
	Convey("Calling DeleteUser with an unauthenticated context must fail", t, func() {
		handler := NewAPIHandler()
		err := handler.DeleteUser(context.Background(), "loki")
		So(err, ShouldEqual, Forbidden)
	})
}

func TestBaseHandler_DeleteUser3(t *testing.T) {
	Convey("Calling DeleteUser with an authenticated context on another user data must fail", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		err := handler.DeleteUser(authenticatedContext, "hulk")
		So(err, ShouldEqual, NotAuthorized)
	})
}

func TestBaseHandler_DeleteUser4(t *testing.T) {
	Convey("Calling DeleteUser with an authenticated context on owned data must succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "loki")
		err := handler.DeleteUser(authenticatedContext, "loki")
		So(err, ShouldBeNil)
	})
}

func TestBaseHandler_DeleteUser5(t *testing.T) {
	Convey("Calling DeleteUser with an admin context must succeed", t, func() {
		handler := NewAPIHandler()
		authenticatedContext := ContextWithLoggedUser(context.Background(), "admin")
		err := handler.DeleteUser(authenticatedContext, "hulk")
		So(err, ShouldBeNil)
	})
}
