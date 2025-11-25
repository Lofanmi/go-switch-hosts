package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ying32/govcl/vcl"
)

var (
	configSwitchHostsDir string
	systemHostsFilename  string

	formMain      *TFormMain
	configManager *ConfigManager
)

func init() {
	configSwitchHostsDir = filepath.Join(GetHomeDir(), ".SwitchHosts")
	systemHostsFilename = EtcHostsFilename()

	if v := os.Getenv("GOSH_SWITCHHOSTSDIR"); v != "" {
		configSwitchHostsDir = v
	}
	if v := os.Getenv("GOSH_HOSTSFILENAME"); v != "" {
		systemHostsFilename = v
	}
}

func main() {
	configManager = NewConfigManager(configSwitchHostsDir)
	if err := configManager.Load(); err != nil {
		vcl.ShowMessage(fmt.Sprintf("加载 SwitchHosts 配置失败: %v", err))
		return
	}
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.SetShowMainForm(true)
	vcl.Application.SetTitle("GoSwitchHosts v1.0")
	vcl.Application.CreateForm(&formMain)
	vcl.Application.SetScaled(true)
	vcl.Application.Run()
}
