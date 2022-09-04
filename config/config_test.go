package config

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestLoad(t *testing.T){
	if _, err:=os.Stat("/home/nikhil/kindleConfig");err!=nil{
		fmt.Println("Config doesn't exists at path , please create one now")
		return
	}
	_, err:=Load("/home/nikhil/kindleConfig")
	fmt.Println("io.EOF ", io.EOF)
	fmt.Println("io.ErrClosedPipe", io.ErrClosedPipe)
	fmt.Println("io.ErrNoProgress", io.ErrNoProgress)
	fmt.Println("io.ShortBuffer", io.ErrShortBuffer)
	fmt.Println("io.ErrShortWrite", io.ErrShortWrite)
	fmt.Println("io.ErrClosedPipe", io.ErrUnexpectedEOF)
	if err!=nil{
		fmt.Println(err)
		t.Fatal(err)
	}
}
