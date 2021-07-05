package ditt

import (
	"bufio"
	"github.com/tidwall/gjson"
	"io"
	"strings"
)

func newJsonObjectStreamParser(reader io.Reader) *jsonObjectStreamParser {
	return &jsonObjectStreamParser{
		reader:   bufio.NewReader(reader),
		userChan: make(chan UserData, 1),
	}
}

type jsonObjectStreamParser struct {
	startConsuming bool
	reader         *bufio.Reader
	userChan       chan UserData
}

func (p *jsonObjectStreamParser) parseUsers(callback UserDataCallback) error {
	var (
		err          error
		parsingUsers bool
		c            rune
		user         UserData
	)
	for {
		c, _, err = p.reader.ReadRune()
		if err != nil {
			return err
		}

		switch c {
		case '[', ',':
			if c == '[' {
				if parsingUsers {
					return BadInput
				}
				parsingUsers = true
			}

			user, err = p.parseUser()
			if err != nil {
				return err
			}

			if !gjson.Valid(string(user)) {
				return BadInput
			}

			err = callback(user)
			if err != nil {
				return err
			}
		case ']':
			if !parsingUsers {
				return BadInput
			}
			return nil
		}
	}
}

func (p *jsonObjectStreamParser) parseUser() (UserData, error) {
	var user string
	metParenthesisAtLeastOnce := false

	expectedClosingBraceCount := 0

	for {
		c, _, err := p.reader.ReadRune()
		if err != nil {
			if err == io.EOF && expectedClosingBraceCount > 0 {
				return "", BadInput
			}
			return "", err
		}

		user = user + string(c)
		switch c {
		case '{':
			expectedClosingBraceCount++
			metParenthesisAtLeastOnce = true
		case '}':
			expectedClosingBraceCount--
		}

		if metParenthesisAtLeastOnce && len(user) > 0 && expectedClosingBraceCount == 0 {
			break
		}
	}
	return UserData(strings.Trim(user, " ")), nil
}
