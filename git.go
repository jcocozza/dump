package main

import (
	"os"
	"os/exec"
)

// generate a git command
//
// always run in root directory as specified in config.go
func gitCmd(cmd string, args ...string) *exec.Cmd {
	argList := append([]string{"-C", root, cmd}, args...)
	return exec.Command("git", argList...)
}

// return a log file to write git stuff to
//
// specified in config.go
//
// MAKE SURE TO CLOSE IT
func gitCmdLog() (*os.File, error) {
	return os.OpenFile(gitCmdLogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
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
	}
}
