package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// generate a git command
//
// always run in root directory as specified in config.go
func gitCmd(cmd string, args ...string) *exec.Cmd {
	argList := append([]string{"-C", config.Root(), cmd}, args...)
	return exec.Command("git", argList...)
}

// return a log file to write git stuff to
//
// specified in config.go
//
// MAKE SURE TO CLOSE IT
func gitCmdLog() (*os.File, error) {
	log := config.GitCmdLog()
	if _, err := os.Stat(log); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(log), 0700)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(config.GitCmdLog(), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
}

func createGitIgnore() error {
	s := fmt.Sprintf("%s\n%s", config.gitCmdLog, ".gitignore")
	contents := []byte(s)
	return os.WriteFile(config.Path(".gitignore"), contents, 0660)
}

func runGit(cmd string, args ...string) error {
	log, err := gitCmdLog()
	if err != nil {
		return err
	}
	defer log.Close()
	c := gitCmd(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = log
	c.Stderr = log
	return c.Run()
}

func insideGit() bool {
	cmd := gitCmd("rev-parse", "--is-inside-work-tree")
	cmd.Stdin = os.Stdin
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func init() {
	in := insideGit()
	if !in {
		err := runGit("init")
		if err != nil {
			panic(err)
		}
		err = createGitIgnore()
		if err != nil {
			panic(err)
		}
	}
}
