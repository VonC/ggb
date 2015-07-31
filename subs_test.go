package main

import "testing"

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
}
