package main

import(
	"fmt"
	// "os"
)


func main() {
	fmt.Printf("Hello\n")
}

func display(token, kind string) {
        fmt.Printf("%-15s %s\n", token, kind)
}