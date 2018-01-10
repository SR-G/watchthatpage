package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PrepareCommand : replace the variables inside the configured commands
func PrepareCommand(command string, cache string, url string, filename string) []string {
	command = strings.Replace(command, "${cache}", cache, -1)
	command = strings.Replace(command, "${url}", url, -1)
	command = strings.Replace(command, "${filename}", filename, -1)
	return strings.Split(command, " ")
}

// GenerateScreenshot : create a screenshot
func GenerateScreenshot(command string, cache string, url string, filename string) {
	replaced := PrepareCommand(command, cache, url, filename)
	binary := replaced[:1][0]
	if _, err := os.Stat(binary); err == nil {
		cmd := exec.Command(binary, replaced[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		fmt.Println("Binary [" + binary + "] not found, can't generate screenshot")
	}
}
