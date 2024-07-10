package main

import (
	"fmt"
	"net/url"
	"strings"
)

// URI represents MongoDB connection URI.
type URI struct {
	Host     string
	Username string
	Password string
	Database string
}

// ParseURI given uri and returns MongoDB connection struct.
func ParseURI(given string) (*URI, error) {
	u, err := url.Parse(given)

	if err != nil {
		return nil, err
	}

	if u.Scheme != "mongodb" && u.Scheme != "mongodb+srv" {
		return nil, fmt.Errorf("scheme must be mongodb")
	}

	pw, exists := u.User.Password()

	if !exists {
		return nil, fmt.Errorf("uri does not contain password")
	}

	paths := strings.Split(u.Path, "/")

	if len(paths) != 2 || paths[1] == "" {
		return nil, fmt.Errorf("incorrect database name")
	}

	return &URI{
		Host:     u.Host,
		Username: u.User.Username(),
		Password: pw,
		Database: paths[1],
	}, nil
}
