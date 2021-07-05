package ditt

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsing(t *testing.T) {
	Convey("Test parsing", t, func() {
		// parsing well formatted json file should end up with parsing ok
		user0 := `{"name": "hulk"}`
		user1 := `{"name": "Loki"}`
		user2 := `{"name": "spider man"}`
		user3 := `{"name": "madara"`

		dataset := fmt.Sprintf("[%s, %s, %s]", user0, user1, user2)

		reader := bytes.NewBufferString(dataset)
		parser := newJsonObjectStreamParser(reader)

		var parsed []UserData
		err := parser.parseUsers(func(data UserData) error {
			parsed = append(parsed, data)
			return nil
		})

		So(err, ShouldBeNil)
		So(parsed, ShouldHaveLength, 3)
		So(user0, ShouldEqual, parsed[0])
		So(user1, ShouldEqual, parsed[1])
		So(user2, ShouldEqual, parsed[2])

		// parsing a json object with a missing closing brace should fail
		dataset = fmt.Sprintf("[%s]", user3)
		reader = bytes.NewBufferString(dataset)
		parser = newJsonObjectStreamParser(reader)
		parsed = nil

		err = parser.parseUsers(func(data UserData) error {
			parsed = append(parsed, data)
			return nil
		})
		So(err, ShouldEqual, BadInput)
	})
}
