package main

import (
	"4ctf/cmd"
	"log"
)

//go:generate sqlboiler --wipe --add-soft-deletes --add-global-variants mysql

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
