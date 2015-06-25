package prj

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var project *Project
var Debug bool

type Project struct {
	rootFolder string
	// Global  GOPATH
	ggopath string
}

func GetProject() (*Project, error) {
	if project == nil {
		project = &Project{}
		gdir, gerr := Git("rev-parse --git-dir")
		gdir = strings.TrimSpace(gdir)
		// fmt.Printf("ko '%s' '%s'", gdir, gerr)
		if gerr != nil {
			return nil, gerr
		}
		if gdir != ".git" {
			project.rootFolder = gdir[:len(gdir)-5]
		} else {
			// fmt.Printf("ok")
			project.rootFolder = wd
		}
		project.ggopath = os.Getenv("GOPATH")
	}
	// fmt.Printf("prf '%s'", project.rootFolder)
	// fmt.Printf("prf '%s'", project.ggopath)
	return project, nil
}

func (p *Project) RootFolder() string {
	return p.rootFolder
}

// Inspired by https://github.com/ghthor/journal/blob/0bd4968a4f9841befdd0dde96b2096e6c930e74c/git/git.go

var gitPath string
var goPath string
var wd string

func init() {
	gitPath = getPathForExe("git")
	goPath = getPathForExe("go")
	var err error
	wd, err = os.Getwd()
	if err != nil {
		log.Fatal("Working directory not accessible")
	}
}

func getPathForExe(exe string) string {
	var err error
	var path = ""
	if path, err = exec.LookPath(exe); err != nil {
		aliases := ""
		if runtime.GOOS == "windows" {
			aliases, err = execcmd("doskey", "/macros")
		} else {
			aliases, err = execcmd("alias", "")
		}
		r := regexp.MustCompile(`(?m)^` + exe + `=(.*)\s+[\$%@\*].*$`)
		sm := r.FindAllStringSubmatch(aliases, 1)
		if len(sm) != 1 || len(sm[0]) != 2 {
			log.Fatalf("Unable to find '%s' path in aliases '%s'", exe)
		}
		return sm[0][1]
	}
	if runtime.GOOS == "windows" {
		if strings.HasSuffix(path, ".bat") {
			bat, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Unable to read '%s' for '%s'", path, exe)
			}
			bats := string(bat)
			r := regexp.MustCompile(`(?m)^\s*?(.*)\s+[\$%@\*].*$`)
			sm := r.FindAllStringSubmatch(bats, 1)
			if len(sm) != 1 || len(sm[0]) != 2 {
				log.Fatalf("Unable to find '%s' path in file '%s'", exe, path)
			}
			return sm[0][1]
		}
	}
	if path == "" {
		log.Fatalf("Unable to get path for '%s'", exe)
	}
	return path
}

// Construct an *exec.Cmd for `git {args}` with a workingDirectory
func Git(cmd string) (string, error) {
	return execcmd(gitPath, cmd)
}
func Golang(cmd string) (string, error) {
	os.Setenv("GOPATH", project.rootFolder+`/deps`)
	os.Setenv("GOBIN", project.rootFolder+`/bin`)
	return execcmd(goPath, cmd)
}

func execcmd(exe, cmd string) (string, error) {
	if Debug {
		fmt.Printf("%s %s\n", exe, cmd)
	}
	args := strings.Split(cmd, " ")
	c := exec.Command(exe, args...)
	c.Dir = project.rootFolder
	var bout bytes.Buffer
	c.Stdout = &bout
	var berr bytes.Buffer
	c.Stderr = &berr
	err := c.Run()
	if err != nil {
		return bout.String(), fmt.Errorf("Unable to run '%s %s' in '%s': err '%s'\n'%s'", exe, cmd, wd, err.Error(), berr.String())
	} else if berr.String() != "" {
		return bout.String(), fmt.Errorf("Warning on run '%s %s' in '%s': '%s'", exe, cmd, wd, berr.String())
	}
	return bout.String(), nil
}
