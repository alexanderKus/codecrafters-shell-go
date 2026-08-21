package main

import (
	"strings"
	"fmt"
	"os"
	"bufio"
)

func echo(args ...string) {
	fmt.Println(strings.Join(args, " "))
}


var run = true
var lookup = map[string]func(args ...string){
	"echo": echo,
	"exit": func(args ...string) {
		run = false
	},
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
