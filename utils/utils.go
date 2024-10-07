package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const version string = "v0.01"

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

func HelpMessage() {
	fmt.Println(Underline + Bold + Blue + "Usage: wiki-go <search_term>" + Reset)
	fmt.Println(Green + "    -h/--help" + Reset + "       Shows this help message")
	fmt.Println(Green + "    -v/--version" + Reset + "    Shows version information")
	fmt.Println(Green + "    -s/--search" + Reset + "     Selects the first result that appears")
}

func VersionMessage() {
	fmt.Println(Blue + "wiki-go version " + version + Reset)
}

// Useless aah function
func ResetStyles() {
	fmt.Print(Reset)
}
