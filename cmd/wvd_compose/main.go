package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yinyajiang/gof/ofdrm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wvd_compose <client_id> <client_private_key>")
		fmt.Println("       wvd_compose <client and private file directory>")
		os.Exit(1)
	}
	var clientID string
	var clientPrivateKey string
	if len(os.Args) == 3 {
		clientID = os.Args[1]
		clientPrivateKey = os.Args[2]
	} else {
		files, _ := filepath.Glob(os.Args[1] + "/*id*")
		if len(files) != 1 {
			fmt.Println("Not found client id file")
			os.Exit(1)
		}
		clientID = files[0]
		files, _ = filepath.Glob(os.Args[1] + "/*private*")
		if len(files) != 1 {
			fmt.Println("Not found client private key file")
			os.Exit(1)
		}
		clientPrivateKey = files[0]
	}

	clientIDByte, err := os.ReadFile(clientID)
	if err != nil {
		fmt.Println("Error reading client ID file:", err)
		os.Exit(1)
	}
	clientPrivateKeyByte, err := os.ReadFile(clientPrivateKey)
	if err != nil {
		fmt.Println("Error reading client private key file:", err)
		os.Exit(1)
	}
	wvd, err := ofdrm.ComposeWVD(clientIDByte, clientPrivateKeyByte)
	if err != nil {
		fmt.Println("Error composing WVD:", err)
		os.Exit(1)
	}
	wvdPath := filepath.Join(filepath.Dir(clientID), "wvd_client")
	err = os.WriteFile(wvdPath, wvd, 0644)
	if err != nil {
		fmt.Println("Error writing WVD file:", err)
		os.Exit(1)
	}
	fmt.Println("WVD file written to:", wvdPath)
}
