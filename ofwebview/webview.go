package ofwebview

import (
	"sync"
	"time"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/webviewloader"
)

type OFWebviewConfig struct {
	WebviewConfig  webviewloader.Config
	Title          string
	Width          int
	Height         int
	WebviewWorkDir string
}

type WebView struct {
	loader          *webviewloader.Loader
	config          OFWebviewConfig
	lock            sync.Mutex
	lastLoginResult LoginResult
}

func NewWebView(cfg OFWebviewConfig) *WebView {
	if cfg.Title == "" {
		cfg.Title = "OnlyFans Login"
	}
	if cfg.Width == 0 {
		cfg.Width = 800
	}
	if cfg.Height == 0 {
		cfg.Height = 600
	}
	if cfg.WebviewConfig.WebviewAppWorkDir == "" {
		cfg.WebviewConfig.WebviewAppWorkDir = cfg.WebviewWorkDir
	}

	return &WebView{
		loader: webviewloader.New(cfg.WebviewConfig),
		config: cfg,
	}
}

func (w *WebView) IsEnable() bool {
	return w.loader.HasMustCfg()
}

func (w *WebView) Install(checkUpdate bool) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.loader.InstallEnv(checkUpdate)
}

func (w *WebView) Check(checkUpdate bool) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.loader.CheckEnv(checkUpdate)
}

func (w *WebView) Login() (LoginResult, error) {
	err := w.Check(false)
	if err != nil {
		err = w.Install(false)
		if err != nil {
			return LoginResult{}, err
		}
	}

	hasLockFailed := false
	for {
		if w.lock.TryLock() {
			break
		}
		time.Sleep(time.Second)
		hasLockFailed = true
	}
	defer w.lock.Unlock()

	if hasLockFailed {
		return w.lastLoginResult, nil
	}

	info, err := w.loader.Start(gof.OFPostDomain, webviewloader.WebviewOptions{
		UA:           "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		Title:        w.config.Title,
		Width:        w.config.Width,
		Height:       w.config.Height,
		WaitElements: []string{".m-logout"},
		WaitCookies:  []string{"sess", "auth_id", "fp"}, // fp is xbc
	})
	if err != nil {
		return LoginResult{}, err
	}
	w.lastLoginResult = LoginResult{
		UA:      info.UA,
		Cookies: info.Cookies,
	}
	return w.lastLoginResult, nil
}
