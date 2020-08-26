package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	ps "github.com/mitchellh/go-ps"
)

func main() {

	pid := 5459
	process, err := getProcessRunningStatus(pid)
	if err != nil {
		log.Fatalf("Error is: %s", err)
	}
	log.Println(process.Pid)

	p, err := ps.FindProcess(pid)

	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(-1)
	}

	fmt.Println("Process ID : ", p.Pid())
	fmt.Println("Parent Process ID : ", p.PPid())
	fmt.Println("Process ID binary name : ", p.Executable())
}

// check if the process is actually running
// However, on Unix systems, os.FindProcess always succeeds and returns
// a Process for the given pid...regardless of whether the process exists
// or not.
func getProcessRunningStatus(pid int) (*os.Process, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	//double check if process is running and alive
	//by sending a signal 0
	//NOTE : syscall.Signal is not available in Windows

	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return proc, nil
	}

	if err == syscall.ESRCH {
		return nil, errors.New("process not running")
	} else {
		return nil, err
	}

	// default
	return nil, errors.New("process running but query operation not permitted")
}
