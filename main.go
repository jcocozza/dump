package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const root = "test"

func getEditor() string {
	ed := os.Getenv("EDITOR")
	if ed == "" {
		ed = "vim"
	}
	return ed
}

func runEditor(name string, fName string) error {
	p := filepath.Join(root, fName)
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

func list(dir string, depth int) error {
	items, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, item := range items {
		idnt := strings.Repeat("\t", depth)
		if item.IsDir() {
			fmt.Fprintf(os.Stdout, "%s%s/\n", idnt, item.Name())
			list(filepath.Join(dir, item.Name()), depth+1)
			continue
		}

		info, err := item.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "%s%s - %s\n", idnt, item.Name(), info.ModTime().Format("2006-01-02"))
	}
	return nil
}

func main() {
	ed := getEditor()
	if len(os.Args) < 2 {
		runEditor(ed, emptyName())
		return
	}
	switch os.Args[1] {
	case "ls":
		list(root, 0)
	default:
		fName := os.Args[1]
		runEditor(ed, fName)
	}
}
