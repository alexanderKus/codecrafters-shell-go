package main

import (
	"fmt"
)

func echo() {
	fmt.Println("echo")
}

var run = true
var lookup = map[string]func(){
	"echo": echo,
	"exit": func() {
		run = false
	},
}

func main() {
	var input string

	for run {
		fmt.Print("$ ")
		fmt.Scanln(&input)

		fn, ok := lookup[input]; 
		if ok {
			fn()
		} else {
			msg := fmt.Sprintf("%v: command not found", input)
			fmt.Println(msg)
		}
	}
}
