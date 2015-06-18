package prj

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var project *Project

type Project struct {
	rootFolder string
}

func GetProject() (*Project, error) {
	if project == nil {
		project = &Project{}
		gdir, gerr := git("rev-parse --git-dir")
		gdir = strings.TrimSpace(gdir)
		// fmt.Printf("ko '%s' '%s'", gdir, gerr)
		if gerr != "" {
			return nil, fmt.Errorf("'%s': '%s'", wd, gerr)
		}
		if gdir != ".git" {
			project.rootFolder = gdir[:len(gdir)-5]
		} else {
			// fmt.Printf("ok")
			project.rootFolder = wd
		}
	}
	// fmt.Printf("prf '%s'", project.rootFolder)
	return project, nil
}

func (p *Project) RootFolder() string {
	return p.rootFolder
}

// Inspired by https://github.com/ghthor/journal/blob/0bd4968a4f9841befdd0dde96b2096e6c930e74c/git/git.go

var gitPath string
var wd string

func init() {
	var err error
	gitPath, err = exec.LookPath("git")
	if err != nil {
		log.Fatal("git must be installed")
	}
	wd, err = os.Getwd()
	if err != nil {
		log.Fatal("Working directory not accessible")
	}
}

// Construct an *exec.Cmd for `git {args}` with a workingDirectory
func git(cmd string) (string, string) {
	args := strings.Split(cmd, " ")
	c := exec.Command(gitPath, args...)
	c.Dir = wd
	var bout bytes.Buffer
	c.Stdout = &bout
	var berr bytes.Buffer
	c.Stderr = &berr
	err := c.Run()
	if err != nil {
		log.Fatalf("Unable to run 'git %s': err '%s'", cmd, err.Error())
	}
	return bout.String(), berr.String()
}
