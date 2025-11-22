package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"flag"
)

func getEditor() string {
	ed := os.Getenv("EDITOR")
	if ed == "" {
		ed = "vim"
	}
	return ed
}

func runEditor(name string, fName string) error {
	p := config.Path(fName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(p), 0700)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(name, p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func emptyName() string {
	now := time.Now()
	s := now.Format("2006-01-02:15:04:05")
	return fmt.Sprintf("%s-%s", "dump", s)
}

// ignore hidden files
func shouldIgnore(name string) bool {
	return name[0] == '.'
}

func list(dir string, depth int, pretty bool) error {
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, item := range items {
		idnt := strings.Repeat("\t", depth)
		if shouldIgnore(item.Name()) {
			continue
		}
		if item.IsDir() {
			if pretty {
				fmt.Fprintf(os.Stdout, "%s%s/\n", idnt, item.Name())
			}
			list(filepath.Join(dir, item.Name()), depth+1, pretty)
			continue
		}
		info, err := item.Info()
		if err != nil {
			return err
		}
		if pretty {
			fmt.Fprintf(os.Stdout, "%s%s - %s\n", idnt, item.Name(), info.ModTime().Format("2006-01-02"))
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", filepath.Join(dir, item.Name()))
		}
	}
	return nil
}

func syncPeer(p peer) {
	err := runGit("fetch", p.name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[warning] unable to sync %s. skipping...\n", p.name)
		return
	}
	err = runGit("pull", "--rebase", p.name, p.branch)	
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync failed for %s. please resolve and resync\n", p.name)
		os.Exit(1)
	}
}

func sync() {
	for _, p := range config.peers {
		syncPeer(p)
	}
}

func main() {
	ed := getEditor()
	if len(os.Args) < 2 {
		runEditor(ed, emptyName())
		return
	}


	switch os.Args[1] {
	case "ls":
		lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
		lsSimple := lsCmd.Bool("l", false, "simple print the list")
		lsCmd.Parse(os.Args[2:])
		list(config.Root(), 0, !*lsSimple)
	case "peers":
		gitPeers()
	case "sync":
		syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
		//syncMerge := syncCmd.Bool("m", false, "do a merge instead of rebase")
		syncCmd.Parse(os.Args[2:])
		sync()
	default:
		/* TODO: flags to add:
		 1. -e editor flag
		 2. -m commit message (maybe??)
		*/
		fName := os.Args[1]
		runEditor(ed, fName)
		runGit("add", fName)
		msg := fmt.Sprintf("file: %s", fName)
		runGit("commit", "-m", msg)
	}
}
