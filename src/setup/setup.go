package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("echo", "Called from Go!")

	// Combine stdout and stderr
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output) // => go version go1.3 darwin/amd64
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\t", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString("[Error]")
		os.Stderr.WriteString(fmt.Sprintf("\n--> Error: %s\n", err.Error()))

	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}
