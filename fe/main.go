package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	distDir := flag.String("dist", "", "path to the dist directory")

	flag.Parse()

	if *distDir == "" {
		log.Printf("--dist must be provided")
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(http.Dir(*distDir)))

	log.Printf("Running frontend on 8001, dir: %s", *distDir)

	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}