package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {

	//Run a jar file
	cmd := exec.Command("java", "-jar", "/Users/jmata/git/spring-boot-demo-app-2/target/spring-boot-app-demo-0.0.1-SNAPSHOT.jar")
	//cmd := exec.Command("java", "-version")

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	//run command in background
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with %s\n", err)
	}

	//outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	// Wait for the process to finish or kill it after a timeout (whichever happens first):
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	//Kill the process after 20 seconds
	case <-time.After(25 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
		log.Println("process killed as timeout reached")
	//Listn to done chanel if the program it is terminated by itself before timeout
	case err := <-done:
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					log.Printf("Core Dump: %t", status.CoreDump())
					log.Fatalf("Exit Status: %d", status.ExitStatus())
				}
			} else {
				log.Fatalf("process finished with error = %v", err)
			}

		}
		log.Print("process finished successfully")
	}
}
