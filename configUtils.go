package main

import ( 
	"path/filepath"
	"os"
)

// return the path of p nested within root
func DumpPath(p string) string {
	return filepath.Join(DumpRoot, p)
}

func GetEditor() string {
	if DumpEditor != "ENV" {
		return DumpEditor
	}
	ed := os.Getenv("EDITOR")
	if ed == "" {
		ed = "vim"
	}
	return ed
}

// return a log file to write stuff to
//
// MAKE SURE TO CLOSE IT (defer f.Close())
func DumpLogFile() (*os.File, error) {
	log := DumpPath(DumpLog)
	if _, err := os.Stat(log); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(log), 0700)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(DumpPath(DumpLog), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
}
