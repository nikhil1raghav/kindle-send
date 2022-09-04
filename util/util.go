package util

import "C"
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



func Scanline()string{
	scanner:=bufio.NewScanner(os.Stdin)
	if scanner.Scan(){
		return scanner.Text()
	}
	color.Red("\nInterrupted")
	os.Exit(1)
	return ""
}

//Scan input and trim
func ScanlineTrim() string{
	return strings.TrimSpace(Scanline())
}
