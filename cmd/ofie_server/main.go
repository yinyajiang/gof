package main

import (
	"context"
	"os"

	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
	"github.com/yinyajiang/gof/ofie"
)

func main() {
	addr := ":8199"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	var cfg struct {
		ofie.Config
		OFAuthInfo ofapi.OFAuthInfo
	}
	err := common.FileUnmarshal("config.json", &cfg)
	if err != nil {
		panic(err)
	}
	dl, err := ofie.NewOFIE(cfg.Config)
	if err != nil {
		panic(err)
	}
	if !cfg.OFAuthInfo.IsEmpty() {
		err = dl.Auth(cfg.OFAuthInfo, true)
		if err != nil {
			panic(err)
		}
	}
	err = dl.CheckAuth()
	if err != nil {
		err = dl.AuthByWebview(true)
		if err != nil {
			panic(err)
		}
	}
	dl.Serve(context.Background(), addr)
}
