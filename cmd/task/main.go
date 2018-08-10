// +build windows

package main

import (
	"fmt"
	"github.com/mpcjanssen/orghandler/pkg/wsl"
	"os"
)

func main() {
	err := wsl.Call("task", os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
