package main

import (
	ascii "asciiart/functions"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	//  Check for correct number of arguments
	banner := ""
	if len(os.Args) < 2 {
		log.Fatalln("Not enough arguments")
	} else if len(os.Args) == 2 {
		banner = "standard"
	} else if len(os.Args) > 3 {
		log.Fatalln("Too many arguments")
	} else {
		banner = os.Args[2]
	}
	textStr := os.Args[1]

	// Check if argument contains only ascii printable characters.
	for _, char := range textStr {
		if !(char >= 32 && char <= 126) {
			fmt.Println("Please only use ascii characters")
			os.Exit(1)
		}
	}

	// Selects correct banner style according to user input argument.
	bannerCall := ""
	switch banner {
	case "standard":
		bannerCall = "standard.txt"
	case "shadow":
		bannerCall = "shadow.txt"
	case "thinkertoy":
		bannerCall = "thinkertoy.txt"
	}

	// Open and defer-close ascii-art file.
	file, err := os.Open("banners/" + bannerCall)
	if err != nil {
		os.Exit(2)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("That closing of file did not happen")
		}
	}()

	// Read file by passing io.Reader interface to buffer Scanner, then passing it to Scan() function,
	// then adds buffer lines outlined by '/n' which are read into 'rows' by Scanner.Text()
	Scanner := bufio.NewScanner(file)
	var rows []string
	for Scanner.Scan() { // Scanner.Scan reads file and splits elements with empty string according to newling (\n) token
		rows = append(rows, Scanner.Text()) // Scanner.text() takes the buffer content and appends string element into string slice 'rows'
	}

	// Colonizes  asciiChars map with key value pairs, where key is decimal value of ascii-art
	// and each element/value is one row of an ascii-art character.
	asciiChrs := make(map[int][]string)
	ascii.MapColonizer(asciiChrs, rows)

	// Handles double '\n' characters and single '\n' character
	outputString := ""
	newlineStr := strings.Split(textStr, "\\n")
	for i, val := range newlineStr {
		if val != "" {
			outputString += ascii.PrintAsciiArt(val, asciiChrs)
		}
		if i < len(newlineStr)-1 {
			outputString += string('\n')
		}
	}
	fmt.Printf("%v", outputString)
}
