package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("test")
	os.Exit(0) // want "os.Exit call in main"
}
