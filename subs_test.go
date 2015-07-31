package main

import (
	"strings"
	"testing"
)

type test struct {
	name string
	arg  string
	err  string
}

func TestAddSub(t *testing.T) {
	// Only github.com/a/b is authorized.
	// project must be accessible.
	// Index must be clean.
	// If already there, check if url is the same:
	// - if not, print warning,
	// -if it is, check if checked out:
	//    - if it is, noop, return err nil.
	//    - if it is not, submodule update for that one only
	// If not already there, test:
	// - successfull add
	// - unsuccessfull add:
	//   - test if added to .gitmodule at least
	// Check link on gopath
	// Check sub's subs:
	//   - if sub's sub is already declared as a sub, test SHA1:
	//     - if same SHA1, no-op
	//     - if differ, print warning
	tests := []*test{
		&test{name: "Only github.com/a/b is authorized.", arg: "a/b/c", err: "doesn't match github.com"},
	}
	var err error
	for _, test := range tests {
		err = addsub(test.arg)
		if err == nil || strings.Contains(err.Error(), test.err) == false {
			t.Errorf("Err '%v', expected '%s'", err, test.err)
		}
	}
}
