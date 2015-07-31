package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/VonC/ggb/prj"
)

type dep struct {
	url  string
	path string
}

func addsub(arg string) error {
	fmt.Printf("addsub %v\n", arg)
	var d *dep = nil
	var err error
	if d, err = newDep(arg); err != nil {
		return err
	}
	fmt.Printf("dep='%+v'\n", d)
	var p *prj.Project
	if p, err = prj.GetProject(); err != nil {
		return err
	}
	fmt.Printf("Project '%s'\n", p.RootFolder())
	return nil
}

var githubre = regexp.MustCompile(`(?m)^github.com/([^/]+)/[^/]+$`)

func newDep(arg string) (*dep, error) {
	if githubre.MatchString(arg) == false {
		return nil, fmt.Errorf("Dependency '%s' doesn't match github.com/xxx/yyy", arg)
	}
	res := githubre.FindStringSubmatch(arg)
	path := "deps/src/" + arg
	path = filepath.FromSlash(path)
	author := res[1]
	url := fmt.Sprintf("https://%s@%s", author, arg)
	d := &dep{path: path, url: url}
	return d, nil
}
