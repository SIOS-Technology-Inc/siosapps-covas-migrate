package main

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	// OK
	func() {
		got, err := parseCommand("./examples/000000001_users.json", "demo")

		if err != nil || got == nil {
			t.Errorf("should pass, error %s", err)
			return
		}

		if got.Version == "" {
			t.Errorf("must pass")
			return
		}

		if got.Admin == "" || got.General == "" {
			t.Errorf("must pass")
			return
		}
	}()
}
