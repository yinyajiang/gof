package ofdrm

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	widevine "github.com/iyear/gowidevine"
	"github.com/yinyajiang/gof/common"
	"golang.org/x/exp/rand"
)

type wvdSt struct {
	_clientIDByte         []byte
	_clientPrivateKeyByte []byte
	_wvdByte              []byte
}

func newWVDFromURI(wvdURI any) (*wvdSt, error) {
	wvd, err := common.ReadURI(wvdURI)
	if err != nil {
		return nil, err
	}
	if len(wvd) == 0 {
		return nil, errors.New("wvd is empty")
	}
	return _newWVD(wvd), nil
}

func newWVDFromRawURI(clientIDURI, clientPrivateKeyURI any) (*wvdSt, error) {
	var clientID, clientPrivateKey []byte

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		clientID, _ = common.ReadURI(clientIDURI)
	}()
	go func() {
		defer wg.Done()
		clientPrivateKey, _ = common.ReadURI(clientPrivateKeyURI)
	}()
	wg.Wait()

	if len(clientID) == 0 || len(clientPrivateKey) == 0 {
		return nil, errors.New("clientID or clientPrivateKey is invalid")
	}
	return _newWVDFromRaw(clientID, clientPrivateKey), nil
}

func newWVDFromZipURI(wvdzipURI, wvdzipMD5URI any, cacheDir string) (*wvdSt, error) {
	localPath := filepath.Join(cacheDir, "wvd.zip")
	localMD5Path := filepath.Join(cacheDir, "wvd.zip.md5")
	latestMD5 := ""

	if fileutil.IsExist(localPath) {
		localMD5, err := fileutil.ReadFileToString(localMD5Path)
		if err == nil {
			latestMD5, err = common.ReadURIString(wvdzipMD5URI)
			if err == nil {
				if strings.EqualFold(localMD5, string(latestMD5)) {
					filesize, _ := fileutil.FileSize(localPath)
					file, _ := os.Open(localPath)
					wvd, err := _selectWVDFromZip(file, filesize)
					if err == nil {
						return wvd, nil
					}
				}
			}
		}
	}
	zipdata, err := common.ReadURI(wvdzipURI)
	if err != nil {
		return nil, err
	}
	if latestMD5 == "" {
		latestMD5, _ = common.ReadURIString(wvdzipMD5URI)
	}
	wvd, err := _selectWVDFromZip(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		return nil, err
	}
	if latestMD5 != "" {
		common.WriteFile(localMD5Path, []byte(latestMD5))
	}
	common.WriteFile(localPath, zipdata)
	return wvd, nil
}

func _selectWVDFromZip(r io.ReaderAt, size int64) (*wvdSt, error) {
	if r == nil || size <= 0 {
		return nil, errors.New("reader or size is invalid")
	}
	zipFile, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	files := []*zip.File{}
	for _, file := range zipFile.File {
		if file.FileInfo().IsDir() {
			continue
		}
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, errors.New("wvd not found in zip")
	}
	file := files[rand.Intn(len(files))]
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	wvd, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return _newWVD(wvd), nil
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

func _loadWVDFromCache(cacheDir string) (*wvdSt, error) {
	if fileutil.IsExist(filepath.Join(cacheDir, "wvd")) {
		wvd, err := os.ReadFile(filepath.Join(cacheDir, "wvd"))
		if err == nil {
			return _newWVD(wvd), nil
		}
	}
	return nil, errors.New("wvd not found")
}

func _newWVDFromRaw(clientID, clientPrivateKey []byte) *wvdSt {
	return &wvdSt{
		_clientIDByte:         clientID,
		_clientPrivateKeyByte: clientPrivateKey,
	}
}

func _newWVD(wvd []byte) *wvdSt {
	return &wvdSt{
		_wvdByte: wvd,
	}
}

func (w *wvdSt) composeWVD() []byte {
	wvd, err := ComposeWVD(w._clientIDByte, w._clientPrivateKeyByte)
	if err != nil {
		return nil
	}
	return wvd
}

func loadWVD(cfg DRMWVDOption) (wvd *wvdSt, err error) {
	defer func() {
		if wvd != nil && len(wvd.WVD()) > 0 {
			wvd.cache(cfg.WVDCacheDir)
		} else {
			wvd, err = _loadWVDFromCache(cfg.WVDCacheDir)
		}
	}()

	if cfg.WVDURI != nil {
		if !common.IsURI(cfg.WVDURI, "zip") {
			wvd, err = newWVDFromURI(cfg.WVDURI)
		} else {
			wvd, err = newWVDFromZipURI(cfg.WVDURI, cfg.WVDMd5URIIfZip, cfg.WVDCacheDir)
		}
		if err == nil {
			return
		}
	}
	if cfg.ClientIDURI != nil && cfg.ClientPrivateKeyURI != nil {
		wvd, err = newWVDFromRawURI(cfg.ClientIDURI, cfg.ClientPrivateKeyURI)
	}
	return
}

func ComposeWVD(clientIDByte, clientPrivateKeyByte []byte) ([]byte, error) {
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
	binary.BigEndian.PutUint16(privateKeyLen, uint16(len(clientPrivateKeyByte)))
	buf = append(buf, privateKeyLen...)

	// 写入 privateKey
	buf = append(buf, clientPrivateKeyByte...)

	// 写入 clientID 长度 (2字节)
	clientIDLen := make([]byte, 2)
	binary.BigEndian.PutUint16(clientIDLen, uint16(len(clientIDByte)))
	buf = append(buf, clientIDLen...)

	// 写入 clientID
	buf = append(buf, clientIDByte...)

	//check
	_, err := widevine.NewDevice(
		widevine.FromWVD(bytes.NewReader(buf)),
	)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func CheckWVD(wvdByte []byte) error {
	_, err := widevine.NewDevice(
		widevine.FromWVD(bytes.NewReader(wvdByte)),
	)
	return err
}

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}
