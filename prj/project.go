package prj

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	name := project.name()
	fmt.Printf("name '%s'\n", name)
	depsPrjroot := project.rootFolder + "/deps/src/" + name
	// fmt.Printf("prf '%s'\n", depsPrjroot)
	var err error
	if depsPrjroot, err = filepath.Abs(depsPrjroot); err != nil {
		return nil, err
	}
	depsPrjdir := filepath.Dir(depsPrjroot)
	if fi, _ := os.Stat(depsPrjdir); fi == nil {
		if err := os.MkdirAll(depsPrjdir, os.ModeDir); err != nil {
			return nil, err
		}
	}
	// fmt.Println(depsPrjdir, depsPrjroot)
	if fi, _ := os.Stat(depsPrjroot); fi == nil {
		if _, err = execcmd("mklink", fmt.Sprintf("/J %s %s", depsPrjroot, project.rootFolder)); err != nil {
			return nil, err
		}
	}
	if err = project.updateGGopath(name); err != nil {
		return nil, err
	}
	return project, nil
}

// either base or remote -v origin
func (p *Project) name() string {
	// git remote show -n origin
	// (?m)^(?:http(?:s)://)?(([^@]+)@)?(.*?)(?:.git)?$
	origin := p.origin()
	if origin != "" {
		return origin
	}
	return filepath.Base(p.RootFolder())

}

// git config --local --get remote.origin.url
// (?m)^\s+Fetch URL: (.*?)$
func (p *Project) origin() string {
	gorg, gerr := Git("config --local --get remote.origin.url")
	// fmt.Printf("gorg='%s', gerr='%+v'", gorg, gerr)
	if gorg == "" || gerr != nil {
		return ""
	}
	r := regexp.MustCompile(`(?m)^(?:http(?:s)://)?(([^@]+)@)?(.*?)(?:.git)?$`)
	sm := r.FindAllStringSubmatch(gorg, 1)
	// fmt.Printf("sm: %+v: %d %d\n", sm, len(sm), len(sm[0]))
	if len(sm) == 1 && len(sm[0]) == 4 {
		return sm[0][3]
	}
	return ""
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
	args = append([]string{"/c", exe}, args...)
	c := exec.Command("cmd", args...)
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

func (p *Project) updateGGopath(path string) error {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	gsrcpdir := p.ggopath + string(filepath.Separator) + "src" + string(filepath.Separator) + dir
	var err error
	if gsrcpdir, err = filepath.Abs(gsrcpdir); err != nil {
		return fmt.Errorf("Unable to get absolute path from '%s'", p.ggopath+string(filepath.Separator)+dir)
	}
	fmt.Println(gsrcpdir)
	if fi, _ := os.Stat(gsrcpdir); fi == nil {
		if err := os.MkdirAll(gsrcpdir, os.ModeDir); err != nil {
			return err
		}
	}
	gsrc := p.ggopath + string(filepath.Separator) + "src" + string(filepath.Separator) + path
	if gsrc, err = filepath.Abs(gsrc); err != nil {
		return fmt.Errorf("Unable to get absolute path from '%s'", p.ggopath+string(os.PathSeparator)+path)
	}
	var fi os.FileInfo
	if fi, _ = os.Stat(gsrc); fi == nil {
		if _, err := execcmd("mklink", fmt.Sprintf("/J %s %s", gsrc, project.rootFolder)); err != nil {
			return err
		}
	} else {
		var l string
		if l, err = execcmd("dir", gsrcpdir); err != nil {
			return fmt.Errorf("Unable to list %s in '%s'", name, gsrcpdir)
		}
		fmt.Println(l)
		r := regexp.MustCompile(fmt.Sprintf(`(?m)<J[UO]NCTION>\s+%s\s+\[([^\]]+)\]\s*$`, name))
		n := r.FindAllStringSubmatch(l, -1)
		fmt.Printf("n='%+v'\n", n)
		if len(n) == 1 {
			pp := n[0][1]
			fmt.Printf("pp='%+v' vs. '%s'\n", pp, p.RootFolder())
			if strings.HasPrefix(pp, p.RootFolder()) {
				return nil
			}
			return fmt.Errorf("'%s' should point to '%s' but points instead to '%s'", gsrc, p.RootFolder(), pp)
		}
		r = regexp.MustCompile(fmt.Sprintf(`(?m)<DIR>\s+%s\s*$`, name))
		n = r.FindAllStringSubmatch(l, -1)
		fmt.Printf("n='%+v'\n", n)
		if len(n) == 1 {
			// move dir
			if err = os.Rename(gsrc, gsrc+".1"); err != nil {
				return fmt.Errorf("Unable to rename '%s'", gsrc)
			}
			// Let's try again, now that the dir has been renamed
			return p.updateGGopath(path)
		} else {
			return fmt.Errorf("Unable to access '%s'", gsrc)
		}
	}
	os.Exit(0)
	return nil
}
