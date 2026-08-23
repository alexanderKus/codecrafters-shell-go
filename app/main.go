package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var PATH = os.Getenv("PATH")
var HOME = os.Getenv("HOME")
var pwd, _ = os.Getwd()
var run = true

var builtin = []string{"echo", "exit", "pwd", "cd", "type"}
var builtInLookup = map[string]func(args ...string){
	"echo": echoFn,
	"exit": exitFn,
	"pwd": pwdFn,
	"cd": cdFn,
	"type": typeFn,
}

func echoFn(args ...string) {
  fmt.Println(strings.Join(args, " "))
}

func exitFn(args ...string) {
	run = false
}

func pwdFn(args ...string) {
	fmt.Println(pwd)
}

func cdFn(args ...string) {
	_exists := func(path string) bool {
		info, err := os.Stat(path)
		return err == nil && info.IsDir()
	}

	commandInput := args[0]
	if !_exists(commandInput) {
		msg := fmt.Sprintf("cd: %s: No such file or directory", commandInput)
		fmt.Println(msg)
		return
	}
	if strings.HasPrefix(commandInput, "/") {
		pwd = commandInput
	} else if strings.HasPrefix(commandInput, "./") {
		temp := filepath.Join(pwd, commandInput[2:])
		if _exists(temp) {
			pwd = temp
		} else {
			msg := fmt.Sprintf("cd: %s: No such file or directory", commandInput)
			fmt.Println(msg)
			return
		}
	} else if strings.HasPrefix(commandInput, "..") { 
		parts :=strings.Split(pwd, "/")
		len := len(parts)
		parent := strings.Join(parts[:len], "/")
		temp := filepath.Join(parent, commandInput[:])
		if _exists(temp) {
			pwd = temp
		} else {
			msg := fmt.Sprintf("cd: %s: No such file or directory", commandInput)
			fmt.Println(msg)
			return
		}
	}

}

func typeFn(args ...string) {
	if slices.Contains(builtin, args[0]) {
		fmt.Println(args[0], "is a shell builtin")
	} else {
		handleNonBuiltin(args...)
	}
}

func handleNonBuiltin(args ...string) {
	commandInput := args[0]
	path, exists := doesCommandExist(commandInput)
	var msg string
	if exists {
		msg = fmt.Sprintf("%v is %v", commandInput, path)
	} else {
		msg = fmt.Sprintf("%v: not found", commandInput)
	}
	fmt.Println(msg)
}

func tryToExec(input ...string) {
	commandInput := input[0]
	path, exists := doesCommandExist(commandInput)
	if exists {
		cmd := exec.Command(path, input[1:]...)
		cmd.Args[0] = commandInput
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		msg := fmt.Sprintf("%v: command not found", commandInput)
		fmt.Println(msg)
	}
}

func doesCommandExist(command string) (path string, exists bool) {
	commandFound, err := getExec(command)
	if err != nil {
		return "", false
	}

	hasPerms, err := hasPermissions(commandFound)
	if err != nil {
		return "", false
	}

	if !hasPerms {
		return "", false
	}

	return commandFound, true
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
	commandPath := filepath.Join(path, command)
	info, err := os.Stat(commandPath)
	if err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0 {
		return commandPath, nil
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
			tryToExec(parts...)
			//msg := fmt.Sprintf("%v: not found", input)
      //fmt.Println(msg)
		}
	}
}
