package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Clear() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Useless aah function
func ResetStyles() {
	fmt.Print(Reset)
}
