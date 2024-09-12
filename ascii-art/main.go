package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var spaceIndexes []int

// printUsage() prints instructions on how to use the program
func printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println(`./ascii-art "Hello!" [stylefile]`)
	fmt.Println("OR")
	fmt.Println(`./ascii-art "Hello!" [stylefile]`)
	fmt.Println("OR")
	fmt.Println(`./ascii-art "Hello!" [stylefile]`)
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println(`go run . "Hello!"`)
	fmt.Println(`go run . "Hello!" thinkertoy`)
	fmt.Println(`go run . --align==center "Hello!" thinkertoy`)
	fmt.Println()
	fmt.Println("If no input is provided, style defaults to 'standard' and align to 'left'.")
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

			// Audit answers demand spaces of width 6, unlike exercise examples where
			// they show a width of 4! So removed this:
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
	// all space banners are the same
	space := "      \n      \n      \n      \n      \n      \n      \n      "
	isSpace := false
	for _, line := range bannerLines {
		rowOut := 0
		if len(line) == 0 {
			out += "\n"
			continue
		}
		for i := 0; i < 8; i++ {
			for _, ban := range line {
				if i == 0 && ban == space {
					isSpace = true
				}

				rowBanner := 0
				for _, r := range ban {
					if r == '\n' {
						rowBanner++
					}

					// write to output when output row matches banner row
					if rowOut == rowBanner && r != '\n' {
						out += string(r)
						// write down where the spaces are on the horizontal output
						// for later justification
						if isSpace {
							spaceIndexes = append(spaceIndexes, len(out))
							isSpace = false
						}
					}
				}
			}
			out += "\n"
			rowOut++
		}
	}
	return
}

// getTermWidth returns the width of the terminal in characters
func getTermWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outSplit := strings.Split(string(out), " ")
	width, _ := strconv.Atoi(outSplit[1][:len(outSplit[1])-1])
	return width
}

// justifyAsciis spreads the words from one side to the other
func justifyAsciis(s string, w int) string {
	spaces := make([]string, len(spaceIndexes))
	rows := strings.Split(s, "\n")
	nuRows := []string{}

	for ir, row := range rows {

		// cut the spacestring to a necessary number of pieces
		if ir%8 == 0 {
			adds := ""

			// put all additional spaces in one string
			for i := 0; i < w-len(row); i++ {
				adds += " "
			}

			//fmt.Println("row, adds, row+adds:", len(row), len(adds), len(row)+len(adds), "w, len spaces, len spaceindexes:", w, len(spaces[0]), len(spaceIndexes))

			onThisRow := 0

			for _, v := range spaceIndexes {
				if v >= (ir)*len(row) && v < ((ir)+1)*len(row) {
					onThisRow++
				}
				//fmt.Println(v, ir, onThisRow, ir*len(row))
			}

			for i := 0; i < onThisRow; i++ {
				if i < onThisRow-1 {
					cut := len(adds) / onThisRow
					spaces[i] = adds[:cut+1]
					adds = adds[cut+1:]
				} else {
					spaces[i] = adds
				}
			}
		}

		diff := 0
		for i, v := range spaceIndexes {
			// don't add spaces to empty line
			if len(row) > 1 {
				if v >= (ir/8)*len(row) && v < ((ir/8)+1)*len(row) {
					//fmt.Println("diff:", diff, "v:", v, "spaces i:", len(spaces[i]))
					row = row[:v+diff] + spaces[i] + row[v+diff:]
					diff += len(spaces[i])
				}

			}
		}
		nuRows = append(nuRows, row)
	}
	return strings.Join(nuRows, "\n")
}

// alignLCR justifies the output text left, center or right, and
// cuts off the part that doesn't fit the screen
func alignLCR(s, a string, w int) string {
	rows := strings.Split(s, "\n")
	nuRows := []string{}

	for _, row := range rows {
		add0 := ""
		add1 := ""

		if a == "left" {
			for i := 0; i < w-len(row); i++ {
				add0 += " "
			}
			row = row + add0
		}

		if a == "center" {
			for i := 0; i < (w-len(row))/2; i++ {
				add0 += " "
				add1 += " "
			}
			if (w-len(row))%2 != 0 {
				add1 += " "
			}
			row = add0 + row + add1
		}

		if a == "right" {
			for i := 0; i < w-len(row); i++ {
				add0 += " "
			}
			row = add0 + row
		}

		nuRows = append(nuRows, row)
	}
	return strings.Join(nuRows, "\n")
}

// main prints the input string as banners in the selected ascii art style,
// including line changes.
func main() {
	align := flag.String("align", "", "specify alignment (left, center, right or justify)")
	flag.Parse()

	if *align != "left" && *align != "center" && *align != "right" && *align != "justify" && *align != "" {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]\n\nExample: go run . --align=right something standard")
		os.Exit(1)
	}

	//args := os.Args[1:]
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Provide at least one string argument")
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
		fmt.Println("too many arguments")
		printUsage()
		os.Exit(1)
	}

	terminalWidth := getTermWidth()

	// get the art style as banners
	stylesAsBanners := getStyleBanners(style)

	// get the lines in input as a slices of banners
	bannerLines := putInputToBanners(input, stylesAsBanners)

	// place banners to horizontal lines
	horizontal := getHorizontalLines(bannerLines)

	// justify when asked for
	if *align == "justify" {
		horizontal = justifyAsciis(horizontal, terminalWidth)
	}

	// align the text left, center or right
	horizontal = alignLCR(horizontal, *align, terminalWidth)

	// print the result
	fmt.Print(horizontal)
}
