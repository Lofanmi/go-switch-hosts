package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const ext = ".hosts"

type ConfigLoader struct{}

func (s *ConfigLoader) Path() (path string) {
	if filepath.IsAbs(path) {
		path = filepath.Clean(path)
	}
	return "."
}

func (s *ConfigLoader) Load(path string, hosts contracts.HostsStore) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		//
	})
	go viper.WatchConfig()

	for _, filename := range viper.GetStringSlice("global.hosts") {
		filename = filepath.Join(path, filename+ext)
		data, _ := os.ReadFile(filename)
		hosts.Parse(string(data))
	}
}

func (s *ConfigLoader) Print(hosts contracts.HostsStore) {
	for _, IP := range hosts.IPs() {
		for _, value := range hosts.Query(IP) {
			fmt.Printf("%s %s\n", IP, value)
		}
		fmt.Println("")
	}
}
