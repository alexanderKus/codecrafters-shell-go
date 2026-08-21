package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"slices"
)



var run = true
var builtin = []string{"echo", "exit", "type"}
var lookup = map[string]func(args ...string){
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
		fmt.Println(args[0], "not found")
	}
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

		fn, ok := lookup[command]; 
		if ok {
			fn(args...)
		} else {
			msg := fmt.Sprintf("%v: command not found", input)
			fmt.Println(msg)
		}
	}
}
