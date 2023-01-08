package store

import (
	"os"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// 可以下载 vscode 扩展, 自动语法高亮:
	// https://marketplace.visualstudio.com/items?itemName=tommasov.hosts
	ext = ".hosts"
)

type ConfigLoader struct {
	Parser     contracts.Parser
	OnChangeFn func(event contracts.ChangeEvent)
}

func NewConfigLoader(parser contracts.Parser) contracts.HostsConfigLoader {
	c := &ConfigLoader{Parser: parser}
	c.Load(c.Path())
	return c
}

func (s *ConfigLoader) Path() (path string) {
	path = gotil.Env(gotil.EnvGoSwitchHostsConfigPath, "$HOME"+string(os.PathSeparator)+".go-switch-hosts")
	path = gotil.ParsePath(path)
	log.WithField("path", path).Debug("配置文件路径")
	return
}

func (s *ConfigLoader) Load(path string) {
	configName := gotil.Env(gotil.EnvGoSwitchHostsConfigName, gotil.DefaultConfigName)
	configType := gotil.Env(gotil.EnvGoSwitchHostsConfigType, gotil.DefaultConfigType)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(path)
	log.WithField("config", path+string(os.PathSeparator)+configName+"."+configType).Debug("开始加载配置文件")
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("err", err).Panic("加载配置文件出错")
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.WithField("change", in).Debug("配置文件发生改变")
		if in.Has(fsnotify.Write) || in.Has(fsnotify.Chmod) {
			s.OnChangeFn(contracts.ChangeEvent{
				Type:  contracts.ChangeTypeConfig,
				Event: in,
			})
		}
	})
	viper.WatchConfig()
}

func (s *ConfigLoader) OnChange(fn func(event contracts.ChangeEvent)) {
	s.OnChangeFn = fn
}
