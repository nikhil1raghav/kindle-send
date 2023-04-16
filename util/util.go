package util

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Red = color.New(color.FgRed)
var RedBold = color.New(color.FgRed).Add(color.Bold)
var Cyan = color.New(color.FgCyan)
var CyanBold = color.New(color.FgCyan).Add(color.Bold)
var Green = color.New(color.FgGreen)
var GreenBold = color.New(color.FgGreen).Add(color.Bold)
var Magenta = color.New(color.FgMagenta)

func Scanline() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	color.Red("\nInterrupted")
	os.Exit(1)
	return ""
}

//Scan input and trim
func ScanlineTrim() string {
	return strings.TrimSpace(Scanline())
}

// ExtractLinks extracts links from a file containing urls
func ExtractLinks(filename string) (links []string) {
	links = make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		Red.Println("Error opening link file", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	return
}
