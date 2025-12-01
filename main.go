package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"io"
	"path/filepath"
	"strings"
	"time"
)

func runEditor(name string, fName string) error {
	p := DumpPath(fName)
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

func list(dir string, depth int, pretty bool, fullPath bool) error {
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// this is just a stupid header
	// may remove it
	if !fullPath && pretty && depth == 0 {
		fmt.Fprintln(os.Stdout, "name last_modified")
	}

	if fullPath && pretty && depth == 0 {
		fmt.Fprintf(os.Stdout, "%s\n", dir)
		depth += 1
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
			list(filepath.Join(dir, item.Name()), depth+1, pretty, fullPath)
			continue
		}
		info, err := item.Info()
		if err != nil {
			return err
		}
		if pretty {
			fmt.Fprintf(os.Stdout, "%s%s %s\n", idnt, item.Name(), info.ModTime().Format("2006-01-02"))
		} else {
			if fullPath {
				fmt.Fprintf(os.Stdout, "%s\n", filepath.Join(dir, item.Name()))
			} else {
				fmt.Fprintf(os.Stdout, "%s\n", item.Name())
			}
		}
	}
	return nil
}

func pullPeer(name string, branch string) {
	err := runGit("fetch", name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[warning] unable to pull %s. skipping...\n", name)
		return
	}
	err = runGit("pull", "--rebase", name, branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pull failed for %s. please resolve and repull\n", name)
		os.Exit(1)
	}
}

func pull() {
	for _, p := range DumpPeers {
		pullPeer(p.Name, p.Branch)
	}
}

func doDump(fName string) {
	ed := GetEditor()
	runEditor(ed, fName)
	runGit("add", fName)
	msg := fmt.Sprintf("file: %s", fName)
	runGit("commit", "-m", msg)
}

// returns path added in dump
func addFile(path string, move bool, force bool, prefix string) (string, error) {
	old, err := os.Open(path)
	if err != nil { return "", err }
	defer old.Close()

	name := filepath.Base(path)
	newPath := DumpPath(filepath.Join(prefix, name))

	err = os.MkdirAll(filepath.Dir(newPath), 0700)	
	if err != nil { return "", err }
	
	flags := os.O_RDWR|os.O_CREATE

	if !force {
		flags |= os.O_EXCL
	}

	f, err := os.OpenFile(newPath, flags, 0660)
	if err != nil { return "", err }
	defer f.Close()

	_, err = io.Copy(f, old)
	if err != nil {
		return "", err
	}

	// we don't actually "move"
	// instead, we just remove the file at the old location
	if move {
		err := os.Remove(path)
		if err != nil { return "", err }
	}
	return newPath, nil
}

func add(files []string, move bool, force bool, prefix string) {
	newPaths := make([]string, len(files))
	for i, f := range files {
		newPath, err := addFile(f, move, force, prefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] failed to add file %s: %s\n", f, err.Error())
			os.Exit(1)
		}
		newPaths[i] = newPath
	}

	
	names := make([]string, len(files))
	for i, f := range newPaths {
		runGit("add", f)
		names[i] = filepath.Base(f)
	}
	fList := strings.Join(names, ", ")
	msg := fmt.Sprintf("files: %s", fList)
	runGit("commit", "-m", msg)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintf(os.Stderr, "%s [OPTIONS] [OPT FILENAME]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [OPTIONS] [COMMAND] [OPTIONS] [ARGS]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  ls    list files")
	fmt.Fprintln(os.Stderr, "  peers list peers")
	fmt.Fprintln(os.Stderr, "  pull  pull from peers")
	fmt.Fprintf(os.Stderr,  "  root  print root (use `cd $(%s root)`)\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "  add   add existing files to dump")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "%s <command> -h for more info on a command\n", os.Args[0])
}

func main() {
	flag.Usage = usage
	flag.StringVar(&DumpRoot, "r", DumpRoot, "root of dump file")
	flag.StringVar(&DumpEditor, "e", DumpEditor, "editor to use (use $EDITOR when set to ENV)")
	flag.Parse()

	// args will contain all non-parsed flags
	// so long as we keep flags distinct we don't have to think too hard about parsing
	args := flag.Args()
	if len(args) == 0 {
		doDump(emptyName())
		return
	}

	switch args[0] {
	case "ls":
		lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
		lsSimple := lsCmd.Bool("l", false, "simple print the list")
		lsFull := lsCmd.Bool("f", false, "full path")
		lsCmd.Parse(args[1:])
		list(DumpRoot, 0, !*lsSimple, *lsFull)
	case "peers":
		peersCmd := flag.NewFlagSet("peers", flag.ExitOnError)
		peersCmd.Parse(args[1:])
		gitPeers()
	case "pull":
		pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
		pullCmd.Parse(args[1:])
		pull()
	case "root":
		rootCmd := flag.NewFlagSet("root", flag.ExitOnError)
		rootCmd.Parse(args[1:])
		fmt.Fprintln(os.Stdout, DumpRoot)
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		mv := addCmd.Bool("mv", false, "move the files (instead of copy)")
		p := addCmd.String("p", "", "set a prefix to nest the files (e.g. foo/bar)")
		force := addCmd.Bool("f", false, "force - will overwrite contents")
		addCmd.Parse(args[1:])
		files := addCmd.Args()
		add(files, *mv, *force, *p)
	default:
		/* TODO: flags to add:
		1. -m commit message (maybe??)
		*/
		fName := args[0]
		doDump(fName)
	}
}
