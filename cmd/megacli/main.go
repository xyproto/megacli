// package main is the main package for the MegaCLI program
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/megacli"
	"github.com/xyproto/vt"
)

const (
	versionString = "MegaCLI 1.0.9"

	startMessage = "---=[ MegaCLI ]=---"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Println(versionString)
			return
		case "-h", "--help":
			fmt.Print(usageString)
			return
		}
	}

	// Initialize vt terminal settings
	vt.Init()

	// Prepare a canvas
	c := vt.NewCanvas()
	defer megacli.Cleanup(c)

	// Handle ctrl-c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		megacli.Cleanup(c)
		os.Exit(1)
	}()

	tty, err := vt.NewTTY()
	if err != nil {
		megacli.Cleanup(c)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer tty.Close()
	tty.SetTimeout(10 * time.Millisecond)

	startdirs := []string{".", env.HomeDir(), "/tmp"}
	curdir, err := megacli.MegaCLI(c, tty, startdirs, startMessage)
	if err != nil && err != megacli.ErrExit {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Write the current directory path to stderr at exit, so that shell scripts can use it
	fmt.Fprintln(os.Stderr, curdir)
}
