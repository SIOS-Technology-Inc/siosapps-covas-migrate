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
	ar, _, _, err := jsonparser.Get(got, "adminCommand")

	if err != nil {
		return nil, err
	}

	out.Admin = string(ar)

	gr, _, _, err := jsonparser.Get(got, "command")

	if err != nil {
		return nil, err
	}

	out.General = string(gr)

	return out, nil
}
