package main

import (
	"context"
	"os"

	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofie"
)

func main() {
	addr := ":8199"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	var cfg ofie.Config
	err := common.FileUnmarshal("config.json", &cfg)
	if err != nil {
		panic(err)
	}
	dl, err := ofie.NewOFIE(cfg)
	if err != nil {
		panic(err)
	}
	dl.Serve(context.Background(), addr)
}
