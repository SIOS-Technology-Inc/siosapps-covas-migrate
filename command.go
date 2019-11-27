package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/buger/jsonparser"
)

// Command represents JSON for migration.
type Command struct {
	Version string
	Admin   string
	General string
}

func parseCommand(filepath, dbname string) (*Command, error) {
	if filepath == "" || dbname == "" {
		return nil, fmt.Errorf("invalid input for parse")
	}

	got, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	// Doing magic overwrite.
	got = []byte(strings.ReplaceAll(string(got), "<db>", dbname))

	out := &Command{}

	// Populate version.
	out.Version = filename(filepath)

	// Extract string representations.
	// You cannot use standard json.Unmarshal command, because
	// the attributes are object and unmarshaller returns error.

	admin, err := func() (string, error) {
		val, typ, _, err := jsonparser.Get(got, "adminCommand")

		if err != nil {
			// Skip when schema may not contain adminCommand.
			if typ != jsonparser.NotExist {
				return "", err
			}
		}

		if typ == jsonparser.NotExist {
			return "", nil
		}

		return string(val), nil
	}()

	if err != nil {
		return nil, err
	}

	out.Admin = admin

	general, err := func() (string, error) {
		val, typ, _, err := jsonparser.Get(got, "command")

		if err != nil {
			// Skip when schema may not contain adminCommand.
			if typ != jsonparser.NotExist {
				return "", err
			}
		}

		if typ == jsonparser.NotExist {
			return "", nil
		}

		return string(val), nil
	}()

	if err != nil {
		return nil, err
	}

	out.General = general

	return out, nil
}
