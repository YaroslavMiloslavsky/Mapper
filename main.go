package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"example.com/sitemapper/parse"
)

func main() {
	//Flag for html file
	fileToOpen := flag.String("html", "ex1.html", "the file we want to parse")

	flag.Parse()
	//Here we process the file
	file, fErr := os.Open(*fileToOpen)
	if fErr != nil {
		log.Fatal(fErr)
		panic("File could not be opened")
	}
	defer file.Close()

	links, err := parse.ParseLinks(file)
	if err != nil {
		log.Fatal("Error")
		os.Exit(1)
	}

	fmt.Println(links)

	headers, err := parse.ParseHeaders(file)
	if err != nil {
		log.Fatal("Error")
		os.Exit(1)
	}

	fmt.Println(headers)
}
