// +build windows

package main

import (
	"fmt"
	"github.com/mpcjanssen/orghandler/pkg/wsl"
	"os"
)

func main() {
	err := wsl.Call("/var/lib/gems/2.5.0/gems/taskwarrior-web-1.1.12/bin/task-web", os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
