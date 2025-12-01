package main

// i think this should be $HOME/.local/share/dump
// DumpRoot should be an absolute path
var DumpRoot string = "/tmp/test"
// DumpLog is just a filename
const DumpLog string = ".dump.log"
// the default editor to use
// when set to "ENV" will use the "EDITOR" environment variable
// if "EDITOR" is unset, then this will default to vim
var DumpEditor string = "ENV"

var DumpPeers = []struct{
	// what to call peer in git
	Name string
	// User on the remote machine
	User string	
	// address of the peer
	Addr string
	// location (on disk) of peer repo	
	Path string		
	// branch to sync with
	Branch string	
}{
}
