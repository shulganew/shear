package main

import "os"

func main() {
	os.Exit(0) // want "strait call os.Exit() in the main method prohibited"
	go anotherFunc()
}

func anotherFunc() {
	os.Exit(0)
}
