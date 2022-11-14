package store

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ext = ".hosts"

type ConfigLoader struct{}

func (s *ConfigLoader) Path() (path string) {
	path = gotil.Env(gotil.EnvGoSwitchHostsConfigPath, "$HOME"+string(os.PathSeparator)+".go-switch-hosts")
	if path == "$PWD" {
		path, _ = os.Getwd()
	} else if strings.Contains(path, "$HOME") {
		path = filepath.Join(gotil.GetHomeDir(), path[5:])
	}
	log.WithField("path", path).Debug("配置文件路径")
	return
}

func (s *ConfigLoader) Load(path string, parser contracts.Parser) {
	configName := gotil.Env(gotil.EnvGoSwitchHostsConfigName, "config")
	configType := gotil.Env(gotil.EnvGoSwitchHostsConfigType, "toml")
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(path)
	log.WithField("config", path+string(os.PathSeparator)+configName+"."+configType).Debug("开始加载配置文件")
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("err", err).Panic("加载配置文件出错")
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		//
	})
	go viper.WatchConfig()

	hostsFiles := viper.GetStringSlice("global.hosts")
	hostsDir := viper.GetString("global.hosts_dir")
	aliasMapping := viper.GetStringMapString("alias")
	log.WithField("hostsFiles", hostsFiles).Debug("用户 hosts 文件列表")
	for _, filename := range hostsFiles {
		alias, ok := aliasMapping[filename]
		if !ok {
			alias = filename
		}
		filename = filepath.Join(path, hostsDir, filename+ext)
		log.WithField("filename", filename).Debugf("加载用户 hosts 文件 [%s]", alias)
		data, e := os.ReadFile(filename)
		if e != nil {
			continue
		}
		parser.Comment(">>> " + alias + " - " + filename)
		parser.Parse(string(data))
		parser.Comment("<<< " + alias + " - " + filename)
		parser.EmptyLine()
	}
}

func (s *ConfigLoader) Print(hosts contracts.HostsStore, buf io.Writer) (err error) {
	log.Debug("[开始] 打印 hosts 文件")
	for _, e := range hosts.List() {
		if _, err = buf.Write([]byte(e.String() + "\n")); err != nil {
			return
		}
	}
	log.Debug("[结束] 打印 hosts 文件")
	return
}
