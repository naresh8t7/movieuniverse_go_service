package main

import (
	"fmt"
	cmd "movieuniverse/pkg/movie/cmd/server"
	"log"
	"os"
)

func main() {
	log.Println("main  server")
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
