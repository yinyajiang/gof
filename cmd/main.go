package main

import (
	"fmt"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/ofapi"
)

func main() {
	authInfo := gof.AuthInfo{
		UserID:    "404514599",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		X_BC:      "1b906354646bf00c32b81239e35a0ec243c314f6",
		Cookie:    "sess=673ibueftja03iqp8usc0l6v3f; auth_id=404514599;",
	}
	rules := gof.Rules{
		AppToken:         "33d57ade8c02dbc5a333db99ff9ae26a",
		StaticParam:      "RyY8GpixStP90t68HWIJ8Qzo745n0hy0",
		Prefix:           "30586",
		Suffix:           "67000213",
		ChecksumConstant: 521,
		ChecksumIndexes: []int{
			0, 2, 3, 7, 7, 8, 8, 10, 11, 13, 14, 16, 17, 17, 17, 19, 19, 20, 21, 21, 23, 23, 24, 24, 27, 27, 29, 30, 31, 34, 35, 39,
		},
	}

	ofapi := ofapi.NewOFAPI(ofapi.OFApiConfig{
		AuthInfo: authInfo,
		Rules:    rules,
	})
	user, err := ofapi.GetMeUserInfo()
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// clientID, err := os.ReadFile("/Volumes/1T 移动硬盘/Downloads/device_client_id_blob (1)")
	// if err != nil {
	// 	panic(err)
	// }
	// privateKey, err := os.ReadFile("/Volumes/1T 移动硬盘/Downloads/device_private_key (1)")
	// if err != nil {
	// 	panic(err)
	// }

	// drm := ofdrm.NewOFDRM(
	// 	ofdrm.OFDRMConfig{
	// 		AuthInfo:         authInfo,
	// 		Rules:            rules,
	// 		ClientID:         clientID,
	// 		ClientPrivateKey: privateKey,
	// 		CDRMProjectServer: []string{
	// 			"https://cdrm-project.com/",
	// 		},
	// 	},
	// )
	// videoURL := "https://cdn3.onlyfans.com/dash/files/d/d0/d0d1ebee28857deb265b91189eeef9b0/0hv2s9qx509rmqqdatxkl.mpd,eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiaHR0cHM6XC9cL2NkbjMub25seWZhbnMuY29tXC9kYXNoXC9maWxlc1wvZFwvZDBcL2QwZDFlYmVlMjg4NTdkZWIyNjViOTExODllZWVmOWIwXC8qIiwiQ29uZGl0aW9uIjp7IkRhdGVMZXNzVGhhbiI6eyJBV1M6RXBvY2hUaW1lIjoxNzMwMzM4NDU4fSwiSXBBZGRyZXNzIjp7IkFXUzpTb3VyY2VJcCI6IjE4MC4xNDkuMjM5LjY0XC8zMiJ9fX1dfQ__,WwRazU-PIpk0chl4L25~v47Xzm46dwBC6K6lRmKxemzxnNq6nfGAWJte1j7d6I5-BYRkXwhGvN7is8Q4wYInnuaibxiGKHvbKYSgES66rXRE5P-LqKYj2CbJf4PFmleXr6PNOqhDHgptx7tZ1fC7CeALK5Hmy7agC4s8jDg1o1DRiOyvlHcq-4bh3ab1v~Fp6VO2pgQy2hEiQGKF9Fm1Xhf6zIBZKiL7BZrO6ijsNFDX2G2GkOqtS2lm8lMgCYDhU6MTOYzGjakCGWosgG1BsgbpWyKsIPHZpebKDrsGOLyc0HV0z2DR31Ypu6Zb4k~wHrDuuOYTf99hEfwEiuWZRQ__,K1JM1KV0NHNR73,3542330144,1349895824"
	// lasetModify, err := drm.GetVideoDecryptedKeyAuto(videoURL)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(lasetModify)
	// data, err := drm.GetVideoLastModified(videoURL)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(data)
}
