package main

import "path/filepath"

// TODO: need a make method
// so that we can have sensible defaults
// path = ".local/share/dump"
// brach = "master"
type peer struct {
	// what to call peer in git
	name string
	// user to sync with
	user string
	// address of server
	addr string
	// location of remote repo
	path string
	// branch to sync with
	branch string
}

// TODO: need sensible defaults for these
type DumpConfig struct {
	root      string
	gitCmdLog string
	peers 	  []peer
}

func (c DumpConfig) Path(p string) string {
	return filepath.Join(c.root, p)
}

func (c DumpConfig) Root() string {
	return c.root
}

func (c DumpConfig) GitCmdLog() string {
	return filepath.Join(c.root, c.gitCmdLog)
}
