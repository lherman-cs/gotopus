//+build !test

package main

import (
	"os"
)

func main() {
	os.Exit(Start(os.Args[0], os.Args[1:]...))
}
