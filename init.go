package main

import (
	"os"
	"path/filepath"

	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	log "github.com/sirupsen/logrus"
)

const defaultLogLevel = "debug"

func initLogger() {
	if level, err := log.ParseLevel(gotil.Env(gotil.EnvGoSwitchHostsLogLevel, defaultLogLevel)); err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(level)
	}
}

func backupEtcHosts(dir string) error {
	hosts := gotil.EtcHostsFilename()
	data, err := os.ReadFile(hosts)
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, "etc.hosts.backup")
	if _, err = os.OpenFile(filename, os.O_RDONLY, 0755); err == nil {
		return nil
	}
	log.WithField("filename", filename).Debug("备份 hosts 文件")
	err = os.WriteFile(filename, data, 0755)
	return err
}

func initConfig() {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir = gotil.Env(gotil.EnvGoSwitchHostsConfigPath, filepath.Join(dir, ".go-switch-hosts"))
	if err = os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	if err = backupEtcHosts(dir); err != nil {
		panic(err)
	}
	config := gotil.Env(gotil.EnvGoSwitchHostsConfigName, gotil.DefaultConfigName) + "." +
		gotil.Env(gotil.EnvGoSwitchHostsConfigType, gotil.DefaultConfigType)
	filename := filepath.Join(dir, config)
	_, err = os.ReadFile(filename)
	if _, ok := err.(*os.PathError); ok || err == os.ErrNotExist {
		data := []byte(`[global]
hosts_dir = "$HOME/.go-switch-hosts/hosts"
hosts = [
    # 请按顺序定义好 hosts 文件，优先级由高到低排序。
    # 可以打开目录 hosts_dir 编辑（文件扩展名为 .hosts）：
	# "系统默认",
    # "开发环境",
    # "线上环境",
]
`)
		if err = os.WriteFile(filename, data, 0755); err != nil {
			panic(err)
		}
		hostsDir := filepath.Join(dir, "hosts")
		if err = os.MkdirAll(hostsDir, 0755); err != nil {
			panic(err)
		}
		// 系统默认.hosts
		filename = filepath.Join(hostsDir, "系统默认.hosts")
		data = []byte("127.0.0.1 localhost\n255.255.255.255 broadcasthost\n::1 localhost")
		if err = os.WriteFile(filename, data, 0755); err != nil {
			panic(err)
		}
		// 开发环境.hosts
		filename = filepath.Join(hostsDir, "开发环境.hosts")
		data = []byte("127.0.0.1 www.baidu.com\n127.0.0.1 im.qq.com")
		if err = os.WriteFile(filename, data, 0755); err != nil {
			panic(err)
		}
		// 线上环境.hosts
		filename = filepath.Join(hostsDir, "线上环境.hosts")
		data = []byte("14.215.177.39 www.baidu.com\n119.147.14.34 im.qq.com")
		if err = os.WriteFile(filename, data, 0755); err != nil {
			panic(err)
		}
	}
}
