package main

import (
	"fmt"
	"os"
	"os/exec"
)

// generate a git command
//
// always run in root directory as specified in config.go
func gitCmd(cmd string, args ...string) *exec.Cmd {
	argList := append([]string{"-C", DumpRoot, cmd}, args...)
	return exec.Command("git", argList...)
}

func createGitIgnore() error {
	s := fmt.Sprintf("%s\n%s", DumpLog, ".gitignore")
	contents := []byte(s)
	return os.WriteFile(DumpPath(".gitignore"), contents, 0660)
}

func runGit(cmd string, args ...string) error {
	log, err := DumpLogFile()
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

func addGitRemote(name string, user string, addr string, path string) error {
	return runGit("remote", "add", name, fmt.Sprintf("%s@%s:%s", user, addr, path))
}

func gitPeers() error {
	cmd := gitCmd("remote", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func insideGit() bool {
	cmd := gitCmd("rev-parse", "--is-inside-work-tree")
	cmd.Stdin = os.Stdin
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

// ensure git is all set up
//
// at some point this can be optimized so that some parts of it don't need to be run each time
// e.g. no need to re-add peers each time
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
	for _, peer := range DumpPeers {
		addGitRemote(peer.Name, peer.User, peer.Addr, peer.Path) // we don't care if it already exists
	}
}
