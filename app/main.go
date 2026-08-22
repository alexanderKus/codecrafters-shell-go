package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var PATH = os.Getenv("PATH")
var run = true

var builtin = []string{"echo", "exit", "type"}
var builtInLookup = map[string]func(args ...string){
	"echo": echoFn,
	"exit": func(args ...string) {
		run = false
	},
	"type": typeFn,
}

func echoFn(args ...string) {
  fmt.Println(strings.Join(args, " "))
}

func typeFn(args ...string) {
	if slices.Contains(builtin, args[0]) {
		fmt.Println(args[0], "is a shell builtin")
	} else {
		tryToExec(args...)
	}
}

func tryToExec(input ...string) {
	commandInput := input[0]
	//args := input[1:]
	commandFound, err := getExec(commandInput)
	if err != nil {
		msg := fmt.Sprintf("%v: command not found", commandInput)
		fmt.Println(msg)
		return
	}

	hasPerms, err := hasPermissions(commandFound)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !hasPerms {
		return
	}

	msg := fmt.Sprintf("%v is %v", commandInput, commandFound)
	fmt.Println(msg)
}

func hasPermissions(path string) (hasPerms bool, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().Perm()&0111 != 0, nil
}


func getExec(command string) (commandPath string, err error) {
	paths := strings.Split(PATH, ":")
	for _, path := range paths {
		//fmt.Printf("Searching through %s\n", path)
		if exec, err := getExecFromPath(path, command); err == nil {
			return exec, nil
		}
	}
	return "", fmt.Errorf("command not found")	
}

func getExecFromPath(path string, command string) (outputPath string, err error) {
	if _, err := os.Stat(path); err == nil {
		if entries, err2 := os.ReadDir(path); err2 == nil {
			for _, entry := range entries {
				if entry.Name() == command {
					return filepath.Join(path, command), nil
				}
			}
		} 
	}
	return "", fmt.Errorf("command not found")
}


func main() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)

	for run {
		fmt.Print("$ ")

		if !scanner.Scan() {
			break
		}

		input = scanner.Text()
		parts := strings.Fields(input)

		command := parts[0]
		args := parts[1:]

		fn, ok := builtInLookup[command];
		if ok {
			fn(args...)
		} else {
			msg := fmt.Sprintf("%v: command not found", input)
      fmt.Println(msg)
		}
	}
}