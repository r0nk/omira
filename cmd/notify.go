package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func Notify(stuff string) {
	command := exec.Command("notify")
	stdin, err := command.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprintf(stdin, stuff)
	err = command.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
