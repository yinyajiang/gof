package ofdrm

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/yinyajiang/gof/common"
)

func LoadClient(cacheDir string, clientIDURI, clientPrivateKeyURI string, cachePriority ...bool) (clientID, clientPrivateKey []byte, err error) {
	if len(cachePriority) > 0 && cachePriority[0] {
		clientID, clientPrivateKey, e := loadCachedClient(cacheDir)
		if e == nil {
			return clientID, clientPrivateKey, nil
		}
	}

	if clientIDURI != "" && clientPrivateKeyURI != "" {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			if !strings.HasPrefix(clientIDURI, "http") && fileutil.IsExist(clientIDURI) {
				clientID, _ = os.ReadFile(clientIDURI)
			} else {
				clientID, _ = common.HttpGet(clientIDURI)
			}
		}()
		go func() {
			defer wg.Done()
			if !strings.HasPrefix(clientPrivateKeyURI, "http") && fileutil.IsExist(clientPrivateKeyURI) {
				clientPrivateKey, _ = os.ReadFile(clientPrivateKeyURI)
			} else {
				clientPrivateKey, _ = common.HttpGet(clientPrivateKeyURI)
			}
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
	fileutil.CreateDir(cacheDir)
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
