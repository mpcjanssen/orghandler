package wsl

import (
	"os"
	"os/exec"
	"strings"
)

//Call allows calling of Windows Subsystem for Linux commands via bash
// it takes care to escape the arguments properly.
func Call(cmd string, args []string) error {
	newArgs := make([]string, len(args)+1)
	newArgs[0] = cmd
	for idx, arg := range args[1:] {
		newArgs[idx+1] = "'" + arg + "'"
		// fmt.Printf("%d%s\n", idx+1, newArgs[idx])
	}

	process := exec.Command("bash", "-c", strings.Join(newArgs, " "))
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	process.Stdin = os.Stdin
	err := process.Run()
	return err
}
