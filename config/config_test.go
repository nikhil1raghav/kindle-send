package config

import (
	"fmt"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	if _, err := os.Stat("/home/nikhil/kindleConfig"); err != nil {
		fmt.Println("Config doesn't exists at path , please create one now")
		return
	}
}
