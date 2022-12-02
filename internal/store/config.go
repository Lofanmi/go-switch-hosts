package store

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ChangeType int

type ChangeEvent struct {
	Type  ChangeType
	Event fsnotify.Event
}

const (
	ext = ".hosts"

	ChangeTypeConfig ChangeType = iota
	ChangeTypeHosts
)

type ConfigLoader struct {
	Parser contracts.Parser
	Exit   bool
}

func NewConfigLoader(parser contracts.Parser) contracts.HostsConfigLoader {
	return &ConfigLoader{
		Parser: parser,
		Exit:   false,
	}
}

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

func (s *ConfigLoader) OnChangeEvent(ce ChangeEvent) {
	path := s.Path()
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
		s.Parser.Comment(">>> " + alias + " - " + filename)
		s.Parser.Parse(string(data))
		s.Parser.Comment("<<< " + alias + " - " + filename)
		s.Parser.EmptyLine()
	}
}

func (s *ConfigLoader) WatchHosts(hostsDir string) (err error) {
	for !s.Exit {
		var hostsWatcher *fsnotify.Watcher
		if hostsWatcher, err = fsnotify.NewWatcher(); err != nil {
			return
		}
		hostsDir, _ = filepath.EvalSymlinks(hostsDir)
		eventsWG := new(sync.WaitGroup)
		eventsWG.Add(1)
		go func() {
			defer eventsWG.Done()
			for {
				select {
				case event, ok := <-hostsWatcher.Events:
					if !ok {
						return // close(hostsWatcher.Events)
					}
					s.OnChangeEvent(ChangeEvent{Type: ChangeTypeHosts, Event: event})
				case e, ok := <-hostsWatcher.Errors:
					if ok {
						log.Printf("hostsWatcher error: %v\n", e)
					}
					return
				}
			}
		}()
		_ = hostsWatcher.Add(hostsDir)
		eventsWG.Wait()
		_ = hostsWatcher.Close()
	}
	return
}
