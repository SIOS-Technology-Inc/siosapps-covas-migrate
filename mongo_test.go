package main

import (
	"fmt"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	os.Setenv("URI", "mongodb://root:password@localhost:30000/demo?authSource=admin")

	got := handler()

	if got == nil {
		t.Error("should not return nil")
		return
	}
}

func TestSetup(t *testing.T) {
	// Setup
	func() {
		db = nil

		if err := handler().Collection(migrationCollection).Drop(ctx()); err != nil {
			panic(fmt.Sprintf("failed to drop collection, %s", err))
		}
	}()

	os.Setenv("URI", "mongodb://root:password@localhost:30000/demo?authSource=admin")

	if err := Setup(); err != nil {
		t.Errorf("should not fail, error %s", err)
	}
}

func TestCurrent(t *testing.T) {
	// Setup
	func() {
		db = nil
		os.Setenv("URI", "mongodb://root:password@localhost:30000/demo?authSource=admin")

		if err := handler().Collection(migrationCollection).Drop(ctx()); err != nil {
			panic(fmt.Sprintf("should not fail, error %s", err))
		}

		if err := Setup(); err != nil {
			panic(fmt.Sprintf("should not fail, error %s", err))
		}
	}()

	// First attempt returns init value.
	func() {
		got, err := Current()

		if err != nil {
			t.Errorf("should not fail, error %s", err)
			return
		}

		if got != migrationInitValue {
			t.Errorf("should be identical, got %s", got)
		}
	}()
}

func TestNext(t *testing.T) {
	first := "000000001_users.json"
	second := "000000002_admins.json"

	// OK with init value on existing directory.
	func() {
		got, err := Next("./examples", migrationInitValue)

		if err != nil {
			t.Errorf("should not fail, error %s", err)
			return
		}

		if got == nil {
			t.Errorf("should be identical, got %s", got)
		}

		if got.Version != first {
			t.Errorf("should be identical, got %s", got)
		}

		if got.General == "" {
			t.Errorf("should exist, %#v", got)
		}

		if got.Admin == "" {
			t.Errorf("should exist, %#v", got)
		}
	}()

	// OK with cursor value.
	func() {
		got, err := Next("./examples", first)

		if err != nil {
			t.Errorf("should not fail, error %s", err)
			return
		}

		if got == nil {
			t.Errorf("should be identical, got %s", got)
		}

		if got.Version != second {
			t.Errorf("should be identical, got %s", got)
		}
	}()

	// Fails on non-existing directory.
	func() {
		_, err := Next("./non-existent", migrationInitValue)

		if err == nil {
			t.Errorf("should not fail, error %s", err)
			return
		}
	}()
}

func TestFilename(t *testing.T) {
	type pattern struct {
		sample string
		exp    string
	}

	pats := []pattern{
		{"/some/dir/00000_users.json", "00000_users.json"},
		{"./00000_users.json", "00000_users.json"},
		{"file:///00000_users.json", "00000_users.json"},
	}

	for idx, p := range pats {
		got := filename(p.sample)

		if got != p.exp {
			t.Errorf("case %d expected %s, got %s", idx, p.exp, got)
		}
	}
}

func TestApply(t *testing.T) {
	first := "no-shard/000000001_create_users.json"

	// OK with init value on existing directory.
	func() {
		// Setup
		func() {
			db = nil
			os.Setenv("URI", "mongodb://root:password@localhost:30000/test-apply?authSource=admin&retrywrites=false")

			if err := handler().Collection(migrationCollection).Drop(ctx()); err != nil {
				panic(fmt.Sprintf("should not fail, error %s", err))
			}

			if err := Setup(); err != nil {
				panic(fmt.Sprintf("should not fail, error %s", err))
			}
		}()

		cmd, err := parseCommand("./examples/"+first, handler().Name())

		if err != nil || cmd == nil {
			t.Errorf("should not fail, error %s", err)
			return
		}

		u, _ := ParseURI(os.Getenv("URI"))

		if err := Apply(cmd, u, "develop"); err != nil {
			t.Errorf("should not fail, error %s", err)
		}
	}()

	// Fails on nil input.
	func() {
		os.Setenv("URI", "mongodb://root:password@localhost:30000/test-apply?authSource=admin&retrywrites=false")
		u, _ := ParseURI(os.Getenv("URI"))
		if err := Apply(nil, u, "develop"); err == nil {
			t.Errorf("must pass")
		}
	}()
}
