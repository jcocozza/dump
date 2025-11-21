package main

import "path/filepath"

type DumpConfig struct {
	root      string
	gitCmdLog string
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

var config = DumpConfig{
	root:      "/tmp/test",
	gitCmdLog: ".dump_gitcmd.log",
}
