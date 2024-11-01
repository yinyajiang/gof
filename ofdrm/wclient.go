package ofdrm

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/yinyajiang/gof/common"
)

func loadClient(cacheDir string, clientIDURL, clientPrivateKeyURL string) (clientID, clientPrivateKey []byte, err error) {
	if clientIDURL != "" && clientPrivateKeyURL != "" {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			clientID, _ = common.HttpGet(clientIDURL)
		}()
		go func() {
			defer wg.Done()
			clientPrivateKey, _ = common.HttpGet(clientPrivateKeyURL)
		}()
		wg.Wait()
		if len(clientID) != 0 && len(clientPrivateKey) != 0 {
			cacheClient(cacheDir, clientID, clientPrivateKey)
			return clientID, clientPrivateKey, nil
		}
	}
	return loadCachedClient(cacheDir)
}

func cacheClient(cacheDir string, clientID, clientPrivateKey []byte) {
	os.WriteFile(filepath.Join(cacheDir, "client_id"), clientID, 0644)
	os.WriteFile(filepath.Join(cacheDir, "client_private_key"), clientPrivateKey, 0644)
}

func loadCachedClient(cacheDir string) (clientID, clientPrivateKey []byte, err error) {
	clientID, err = os.ReadFile(filepath.Join(cacheDir, "client_id"))
	if err != nil {
		return nil, nil, err
	}
	clientPrivateKey, err = os.ReadFile(filepath.Join(cacheDir, "client_private_key"))
	if err != nil {
		return nil, nil, err
	}
	return clientID, clientPrivateKey, nil
}
