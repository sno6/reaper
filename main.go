package main

import (
	"log"

	"github.com/sno6/reaper/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
