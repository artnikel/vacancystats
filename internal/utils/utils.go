package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func ClearConsole() {
	cmd := exec.Command("bash", "-c", "clear") // linux
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	arch := false
	if err != nil {
		cmd = exec.Command("cmd", "/c", "cls") // windows
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		arch = true
	}
	if arch {
		err = cmd.Run()
		if err != nil {
			fmt.Printf("\nclear console error: %v\n", err)
		}
	}
}
