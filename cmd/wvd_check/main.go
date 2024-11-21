package main

import (
	"fmt"
	"os"

	"github.com/yinyajiang/gof/ofdrm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wvd_check <wvd_file>")
		os.Exit(1)
	}
	wvdPath := os.Args[1]
	wvdByte, err := os.ReadFile(wvdPath)
	if err != nil {
		fmt.Println("Error reading WVD file:", err)
		os.Exit(1)
	}
	err = ofdrm.CheckWVD(wvdByte)
	if err != nil {
		fmt.Println("Error checking WVD:", err)
		os.Exit(1)
	}
	fmt.Println("WVD file is valid")
}
