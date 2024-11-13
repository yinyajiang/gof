package ofdrm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
	widevine "github.com/iyear/gowidevine"
	"github.com/yinyajiang/gof/common"
)

type wvdSt struct {
	_clientIDByte         []byte
	_clientPrivateKeyByte []byte
	_wvdByte              []byte
}

func newWVDFromRaw(clientID, clientPrivateKey []byte) *wvdSt {
	return &wvdSt{
		_clientIDByte:         clientID,
		_clientPrivateKeyByte: clientPrivateKey,
	}
}

func newWVDFromWVD(wvd []byte) *wvdSt {
	return &wvdSt{
		_wvdByte: wvd,
	}
}

func newWVDFromURI(wvdURI string) (*wvdSt, error) {
	if wvdURI == "" {
		return nil, errors.New("wvdURI cannot be empty")
	}
	var wvd []byte
	if wvdURI != "" {
		if !strings.HasPrefix(wvdURI, "http") && fileutil.IsExist(wvdURI) {
			wvd, _ = os.ReadFile(wvdURI)
		} else {
			wvd, _ = common.HttpGet(wvdURI)
		}
	}
	if len(wvd) == 0 {
		return nil, errors.New("wvd is empty")
	}
	return newWVDFromWVD(wvd), nil
}

func newWVDFromRawURI(clientIDURI, clientPrivateKeyURI string) (*wvdSt, error) {
	if clientIDURI == "" || clientPrivateKeyURI == "" {
		return nil, errors.New("clientIDURI and clientPrivateKeyURI cannot be empty")
	}
	var clientID, clientPrivateKey []byte

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

	if len(clientID) == 0 || len(clientPrivateKey) == 0 {
		return nil, errors.New("clientID or clientPrivateKey is empty")
	}
	return newWVDFromRaw(clientID, clientPrivateKey), nil
}

func newWVDFromCache(cacheDir string) (*wvdSt, error) {
	if fileutil.IsExist(filepath.Join(cacheDir, "client_id")) && fileutil.IsExist(filepath.Join(cacheDir, "client_private_key")) {
		clientID, err := os.ReadFile(filepath.Join(cacheDir, "client_id"))
		if err == nil {
			clientPrivateKey, err := os.ReadFile(filepath.Join(cacheDir, "client_private_key"))
			if err == nil {
				return newWVDFromRaw(clientID, clientPrivateKey), nil
			}
		}
	}
	if fileutil.IsExist(filepath.Join(cacheDir, "wvd")) {
		wvd, err := os.ReadFile(filepath.Join(cacheDir, "wvd"))
		if err == nil {
			return newWVDFromWVD(wvd), nil
		}
	}
	return nil, errors.New("wvd not found")
}

func (w *wvdSt) WVD() []byte {
	if w._wvdByte != nil {
		return w._wvdByte
	}
	w._wvdByte = w.composeWVD()
	return w._wvdByte
}

func (w *wvdSt) cache(cacheDir string) error {
	fileutil.CreateDir(cacheDir)
	return os.WriteFile(filepath.Join(cacheDir, "wvd"), w.WVD(), 0644)
}

func (w *wvdSt) composeWVD() []byte {
	buf := make([]byte, 0)

	type wvdHeader struct {
		Signature     [3]byte
		Version       uint8
		Type          uint8
		SecurityLevel uint8
		Flags         byte
	}
	header := wvdHeader{
		Signature:     [3]byte{'W', 'V', 'D'},
		Version:       2,
		Type:          0, // 默认值
		SecurityLevel: 0, // 默认值
		Flags:         0, // 默认值
	}

	headerBytes := make([]byte, 7)
	copy(headerBytes[0:3], header.Signature[:])
	headerBytes[3] = header.Version
	headerBytes[4] = header.Type
	headerBytes[5] = header.SecurityLevel
	headerBytes[6] = header.Flags
	buf = append(buf, headerBytes...)

	// 写入 privateKey 长度 (2字节)
	privateKeyLen := make([]byte, 2)
	binary.BigEndian.PutUint16(privateKeyLen, uint16(len(w._clientPrivateKeyByte)))
	buf = append(buf, privateKeyLen...)

	// 写入 privateKey
	buf = append(buf, w._clientPrivateKeyByte...)

	// 写入 clientID 长度 (2字节)
	clientIDLen := make([]byte, 2)
	binary.BigEndian.PutUint16(clientIDLen, uint16(len(w._clientIDByte)))
	buf = append(buf, clientIDLen...)

	// 写入 clientID
	buf = append(buf, w._clientIDByte...)

	//check
	_, err := widevine.NewDevice(
		widevine.FromWVD(bytes.NewReader(buf)),
	)
	if err != nil {
		return nil
	}
	return buf
}

func loadWVD(cfg DRMWVDOption) (wvd *wvdSt, err error) {
	if cfg.CachePriority {
		wvd, e := newWVDFromCache(cfg.ClientCacheDir)
		if e == nil {
			return wvd, nil
		}
	}

	save := true
	defer func() {
		if save {
			wvd.cache(cfg.ClientCacheDir)
		}
	}()

	if cfg.WVD != nil {
		return newWVDFromWVD(cfg.WVD), nil
	}
	if cfg.RawWVDID != nil && cfg.RawWVDPrivateKey != nil {
		return newWVDFromRaw(cfg.RawWVDID, cfg.RawWVDPrivateKey), nil
	}

	if cfg.WVDURI != "" {
		wvd, e := newWVDFromURI(cfg.WVDURI)
		if e == nil {
			return wvd, nil
		}
	}
	if cfg.ClientIDURI != "" && cfg.ClientPrivateKeyURI != "" {
		wvd, e := newWVDFromRawURI(cfg.ClientIDURI, cfg.ClientPrivateKeyURI)
		if e == nil {
			return wvd, nil
		}
	}
	save = false
	return newWVDFromCache(cfg.ClientCacheDir)
}
