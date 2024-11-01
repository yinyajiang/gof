package main

import (
	"fmt"
	"os"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/ofapi"
	"github.com/yinyajiang/gof/ofdrm"
)

func main() {
	authInfo := gof.AuthInfo{
		UserID:    "404514599",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		X_BC:      "1b906354646bf00c32b81239e35a0ec243c314f6",
		Cookie:    "sess=673ibueftja03iqp8usc0l6v3f; auth_id=404514599;",
	}

	api, err := ofapi.NewOFAPI(ofapi.Config{
		AuthInfo: authInfo,
	})
	if err != nil {
		panic(err)
	}
	// err := api.CheckAuth()
	// if err != nil {
	// 	panic(err)
	// }

	// err := api.GetPurchased()

	subs, err := api.GetSubscriptions(ofapi.SubscritionTypeActive)
	if err != nil {
		panic(err)
	}
	drmInfo := ofdrm.DRMInfo{}
h:
	for _, sub := range subs {
		posts, err := api.GetUserMedias(sub.ID)
		if err != nil {
			panic(err)
		}
		for _, post := range posts {
			for _, media := range post.Media {
				if media.Type == "video" {
					if media.Files != nil && media.Files.Drm.Manifest.Dash != "" {
						drmInfo = ofdrm.DRMInfo{
							Drm:     media.Files.Drm,
							MediaID: media.ID,
							PostID:  post.ID,
						}
						break h
					}
				}
			}
		}
	}

	clientID, err := os.ReadFile("/Volumes/1T 移动硬盘/Downloads/device_client_id_blob (1)")
	if err != nil {
		panic(err)
	}
	privateKey, err := os.ReadFile("/Volumes/1T 移动硬盘/Downloads/device_private_key (1)")
	if err != nil {
		panic(err)
	}

	drm := ofdrm.NewOFDRM(
		api.Req(),
		ofdrm.OFDRMConfig{
			ClientID:         clientID,
			ClientPrivateKey: privateKey,
		},
	)
	lasetModify, err := drm.GetVideoDecryptedKeyAuto(drmInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(lasetModify)
	data, err := drm.GetVideoLastModified(drmInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
