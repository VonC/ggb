package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/VonC/ggb/prj"
)

type test struct {
	name string
	arg  string
	err  string
}

var currentTest *test

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
	prjGetProject = testPrjGetProject
	prjGit = testPrjGit
	tests := []*test{
		&test{name: "Only github.com/a/b is authorized.",
			arg: "a/b/c", err: "doesn't match github.com"},
		&test{name: "fail get project.",
			arg: "github.com/get/project", err: "Unable to get project"},
		&test{name: "Index must be clean: err",
			arg: "github.com/err/index"},
		&test{name: "Index must be clean: not clean",
			arg: "github.com/notclean/index"},
		&test{name: "Add sub",
			arg: "github.com/add/sub"},
	}
	var err error
	for _, test := range tests {
		currentTest = test
		if currentTest.name == "Add sub" {
			verbose = true
		}
		err = addsub(test.arg)
		if test.err != "" && (err == nil || strings.Contains(err.Error(), test.err) == false) {
			t.Errorf("Err '%v', expected '%s'", err, test.err)
		}
		verbose = false
	}
}

type tproject struct {
	prj.Project
}

func (tp *tproject) RootFolder() string {
	return "testrf"
}

func testPrjGetProject() (prj.Project, error) {
	if currentTest.arg == "github.com/get/project" {
		return nil, fmt.Errorf("Unable to get project")
	}
	return &tproject{}, nil
}

func testPrjGit(cmd string) (string, error) {
	if currentTest.arg == "github.com/err/index" {
		return "", fmt.Errorf("Unable to execute git status")
	}
	if currentTest.arg == "github.com/notclean/index" {
		return "not clean", nil
	}
	return "", nil
}
