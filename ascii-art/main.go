package main

import (
	"fmt"
	"io"
	"os"
)

// printUsage() prints instructions on how to use the program
func printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println(`./ascii-art "Hello!" [stylefile]`)
	fmt.Println("Example:")
	fmt.Println(`go run . "Hello!" thinkertoy`)
	fmt.Println("If no style is provided, it defaults to 'standard'.")
}

// getStyleFile opens and reads the ascii art file based on the provided style argument
func getStyleFile(style string) *os.File {
	file, err := os.Open(style + ".txt")
	if err != nil {
		fmt.Printf("Error opening style file: %s\n", err)
		os.Exit(1)
	}
	return file
}

// getStyleString returns the ascii art style file content as a string
func getStyleString(style string) string {
	styleFile := getStyleFile(style)
	defer styleFile.Close()

	bytes, err := io.ReadAll(styleFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

// getStyleBanners returns the ascii art style as slices of banners
func getStyleBanners(style string) []string {
	styleString := getStyleString(style)
	styleAsBanners := []string{""}

	var prevRu rune
	char := 0
	for j, ru := range styleString {
		if ru == 13 || j == 0 { // skip vertical tabs and the very first line
			continue
		}

		if ru == '\n' && prevRu == '\n' {
			// remove last line change and
			// start a new banner at two line changes
			styleAsBanners[char] = styleAsBanners[char][:len(styleAsBanners[char])-1]
			styleAsBanners = append(styleAsBanners, "")
			char++
		} else {
			// otherwise add rune to current banner
			styleAsBanners[char] += string(ru)
		}
		prevRu = ru
	}
	return styleAsBanners
}

// putInputToBanners retrieves the banners corresponding to the input characters
// and puts them into slices representing lines of text
func putInputToBanners(input string, banners []string) [][]string {
	// multiple slices of banners, separated by line changes
	bannerLines := [][]string{}
	lineIndex := 0
	for i, ru := range input {
		// start a new slice of banners when encountering "\n"
		if ru == 'n' && i > 0 && input[i-1] == '\\' {
			bannerLines[lineIndex] = bannerLines[lineIndex][:len(bannerLines[lineIndex])-1] // remove last banner (`\`)

			if !(len(bannerLines) == 1 && len(bannerLines[lineIndex]) == 0) {
				bannerLines = append(bannerLines, []string{}) // start new line
			}
			lineIndex++
			continue
		}

		if ru > 31 && ru < 127 {
			if len(bannerLines) == 0 {
				bannerLines = append(bannerLines, []string{})
			}

			bannerLines[lineIndex] = append(bannerLines[lineIndex], banners[ru-32])

			/* 			if ru != ' ' {
			   				bannerLines[lineIndex] = append(bannerLines[lineIndex], banners[ru-32])
			   			} else {
			   				// requested spaces are thinner than in styles
			   				bannerLines[lineIndex] = append(bannerLines[lineIndex], "    \n    \n    \n    \n    \n    \n    \n    \n")
			   			} */
		}
	}
	return bannerLines
}

// getHorizontalLines composes horizontal lines from slices of banners
func getHorizontalLines(bannerLines [][]string) (out string) {
	for _, line := range bannerLines {
		rowOut := 0
		if len(line) == 0 {
			out += "\n"
			continue
		}
		for i := 0; i < 8; i++ {
			for _, ban := range line {
				rowBanner := 0
				for _, r := range ban {
					if r == '\n' {
						rowBanner++
					}

					if rowOut == rowBanner && r != '\n' {
						out += string(r)
					}
				}
			}
			out += "\n"
			rowOut++
		}
	}
	return
}

// main prints the input string as banners in the selected ascii art style,
// including line changes.
func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Provide at least a string argument")
		printUsage()
		os.Exit(1)
	}

	input := args[0]
	style := "standard" // Default to 'standard' style

	// Check if a second argument (style) is provided
	if len(args) == 2 {
		style = args[1]
	}

	if len(args) > 2 {
		fmt.Println("to many arguments")
		printUsage()
		os.Exit(1)
	}

	// get the art style as banners
	stylesAsBanners := getStyleBanners(style)

	// get the lines in input as a slices of banners
	bannerLines := putInputToBanners(input, stylesAsBanners)

	// place banners to horizontal lines and print the result
	fmt.Print(getHorizontalLines(bannerLines))
}
