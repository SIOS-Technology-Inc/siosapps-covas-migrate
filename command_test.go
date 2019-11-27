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
			t.Errorf("must pass, %#v", got)
			return
		}
	}()

	// OK when adminCommand does not exist.
	func() {
		got, err := parseCommand("./examples/000000004-admin-only.json", "demo")

		if err != nil || got == nil {
			t.Errorf("should pass, error %s", err)
			return
		}

		if got.Version == "" {
			t.Errorf("must pass")
			return
		}

		if got.Admin == "" || got.General != "" {
			t.Errorf("must pass")
			return
		}
	}()

	// OK when command does not exist.
	func() {
		got, err := parseCommand("./examples/000000005-command-only.json", "demo")

		if err != nil || got == nil {
			t.Errorf("should pass, error %s", err)
			return
		}

		if got.Version == "" {
			t.Errorf("must pass")
			return
		}

		if got.Admin != "" || got.General == "" {
			t.Errorf("must pass")
			return
		}
	}()
}
