package main

import (
	"fmt"
)

func echo() {
	fmt.Println("echo")
}

var lookup = map[string]func(){
	"echo": echo,
}

func main() {
	var input string

	fmt.Print("$ ")
	fmt.Scanln(&input)

	fn, ok := lookup[input]; 
	if ok {
		fn()
	} else {
		msg := fmt.Sprintf("{%v}: command not found", input)
		fmt.Println(msg)
	}
}
