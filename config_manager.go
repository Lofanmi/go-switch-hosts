package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfigEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	On    bool   `json:"on"`
}

type HostsData struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	IDNumber string `json:"_id"`
}

type ConfigManager struct {
	basePath  string
	config    []ConfigEntry
	hostsList []HostsData
	ids       []string
}

func NewConfigManager(basePath string) *ConfigManager {
	return &ConfigManager{
		basePath: basePath,
	}
}

func (cm *ConfigManager) Load() error {
	var config []ConfigEntry
	if err := cm.loadJSONFile(&config, cm.basePath, "data", "list", "tree.json"); err != nil {
		return fmt.Errorf("加载配置列表失败: %w", err)
	}
	cm.config = config

	var ids []string
	if err := cm.loadJSONFile(&ids, cm.basePath, "data", "collection", "hosts", "ids.json"); err != nil {
		return fmt.Errorf("加载hosts ID列表失败: %w", err)
	}
	cm.ids = ids

	hostsList := make([]HostsData, len(ids))
	for i, idNumber := range ids {
		if err := cm.loadJSONFile(&hostsList[i], cm.basePath, "data", "collection", "hosts", "data", idNumber+".json"); err != nil {
			return fmt.Errorf("加载hosts数据 %s 失败: %w", idNumber, err)
		}
	}
	cm.hostsList = hostsList

	return nil
}

func (cm *ConfigManager) SaveToFile(outputFile string) error {
	if outputFile == "" {
		return fmt.Errorf("输出文件路径未设置")
	}

	var enabledHosts strings.Builder
	_, _ = enabledHosts.WriteString(ParseSystemHosts(systemHostsFilename))

	on := false
	for _, config := range cm.config {
		if config.On {
			on = true
			break
		}
	}
	if on {
		_, _ = enabledHosts.WriteString(LineEnding + LineEnding + SwitchHostsContentStart + LineEnding + LineEnding)
		for i, config := range cm.config {
			if !config.On {
				continue
			}
			for _, hostsData := range cm.hostsList {
				if hostsData.ID != config.ID {
					continue
				}
				_, _ = enabledHosts.WriteString("# -------------------------------------------------------------------------------" + LineEnding)
				_, _ = enabledHosts.WriteString(fmt.Sprintf("# %s [%s]%s", config.Title, hostsData.ID, LineEnding))
				_, _ = enabledHosts.WriteString("# -------------------------------------------------------------------------------" + LineEnding)
				_, _ = enabledHosts.WriteString(LineEnding)
				_, _ = enabledHosts.WriteString(strings.TrimSpace(hostsData.Content))
				_, _ = enabledHosts.WriteString(LineEnding)
				if i != len(cm.config)-1 {
					_, _ = enabledHosts.WriteString(LineEnding)
				}
				break
			}
		}
	}
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	if err := os.WriteFile(outputFile, []byte(enabledHosts.String()), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}

func (cm *ConfigManager) GetConfig() []ConfigEntry {
	return cm.config
}

func (cm *ConfigManager) GetHostsByID(id string) (string, bool) {
	for _, hostsData := range cm.hostsList {
		if hostsData.ID == id {
			return hostsData.Content, true
		}
	}
	return "", false
}

func (cm *ConfigManager) UpdateHostsContent(id, content string) error {
	for i, hostsData := range cm.hostsList {
		if hostsData.ID == id {
			cm.hostsList[i].Content = content
			// 保存到文件
			filename := filepath.Join(cm.basePath, "data", "collection", "hosts", "data", hostsData.IDNumber+".json")
			data, err := json.MarshalIndent(cm.hostsList[i], "", "  ")
			if err != nil {
				return fmt.Errorf("序列化数据失败: %w", err)
			}
			if e := os.WriteFile(filename, data, 0644); e != nil {
				return fmt.Errorf("保存文件失败: %w", e)
			}
			return nil
		}
	}
	return fmt.Errorf("未找到ID为 %s 的hosts配置", id)
}

func (cm *ConfigManager) UpdateConfigStatus(id string, on bool) error {
	for i := range cm.config {
		if cm.config[i].ID == id {
			cm.config[i].On = on
			// 保存到配置文件
			filename := filepath.Join(cm.basePath, "data", "list", "tree.json")
			data, err := json.MarshalIndent(cm.config, "", "  ")
			if err != nil {
				return fmt.Errorf("序列化配置失败: %w", err)
			}
			if e := os.WriteFile(filename, data, 0644); e != nil {
				return fmt.Errorf("保存配置文件失败: %w", e)
			}
			return nil
		}
	}
	return fmt.Errorf("未找到ID为 %s 的配置项", id)
}

func (cm *ConfigManager) loadJSONFile(o any, element ...string) error {
	filename := filepath.Join(element...)
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}
