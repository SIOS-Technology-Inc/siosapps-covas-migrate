package main

import (
	"reflect"
	"testing"
)

func TestURI(t *testing.T) {
	// OKs
	func() {
		type pattern struct {
			given string
			exp   *URI
		}

		patterns := []pattern{
			pattern{
				given: "mongodb://user:pass@localhost:27017/db",
				exp: &URI{
					Host:     "localhost:27017",
					Username: "user",
					Password: "pass",
					Database: "db",
				},
			},
		}

		for idx, p := range patterns {
			got, err := ParseURI(p.given)

			if err != nil {
				t.Errorf("case %d should not fail", idx)
				continue
			}

			if !reflect.DeepEqual(got, p.exp) {
				t.Errorf("case %d should be identical\n exp %#v \n got %#v", idx, p.exp, got)
			}
		}
	}()

	// Fails
	func() {
		patterns := []string{
			"http://user:pass@localhost:27017/db",
			"mongodb://user:pass@localhost:27017/",
			"mongodb://user:pass@localhost:27017/db/wrong",
		}

		for idx, p := range patterns {
			got, err := ParseURI(p)

			if err == nil {
				t.Errorf("case %d should fail, got %#v", idx, got)
			}
		}
	}()
}
