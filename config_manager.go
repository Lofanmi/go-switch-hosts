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
			return cm.saveConfigToFile()
		}
	}
	return fmt.Errorf("未找到ID为 %s 的配置项", id)
}

func (cm *ConfigManager) saveConfigToFile() error {
	configFile := filepath.Join(cm.basePath, "data", "list", "tree.json")
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}
	return nil
}

func (cm *ConfigManager) saveIdsToFile() error {
	idsFile := filepath.Join(cm.basePath, "data", "collection", "hosts", "ids.json")
	idsData, err := json.MarshalIndent(cm.ids, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化ID列表失败: %w", err)
	}
	if err = os.WriteFile(idsFile, idsData, 0644); err != nil {
		return fmt.Errorf("保存ID列表文件失败: %w", err)
	}
	return nil
}

func (cm *ConfigManager) AddConfig(config ConfigEntry, content string) error {
	cm.config = append(cm.config, config)
	hostsData := HostsData{
		ID:       config.ID,
		Content:  content,
		IDNumber: config.ID,
	}
	cm.hostsList = append(cm.hostsList, hostsData)
	cm.ids = append(cm.ids, config.ID)
	// 确保目录存在
	if err := os.MkdirAll(filepath.Join(cm.basePath, "data", "list"), 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(cm.basePath, "data", "collection", "hosts", "data"), 0755); err != nil {
		return fmt.Errorf("创建hosts数据目录失败: %w", err)
	}
	// 保存配置文件
	if err := cm.saveConfigToFile(); err != nil {
		return err
	}
	// 保存ID列表文件
	if err := cm.saveIdsToFile(); err != nil {
		return err
	}
	// 保存hosts数据文件
	hostsDataFile := filepath.Join(cm.basePath, "data", "collection", "hosts", "data", config.ID+".json")
	hostsDataContent, err := json.MarshalIndent(hostsData, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化hosts数据失败: %w", err)
	}
	if err = os.WriteFile(hostsDataFile, hostsDataContent, 0644); err != nil {
		return fmt.Errorf("保存hosts数据文件失败: %w", err)
	}
	return nil
}

func (cm *ConfigManager) DeleteConfig(id string) error {
	// 从配置列表中删除
	for i, config := range cm.config {
		if config.ID == id {
			cm.config = append(cm.config[:i], cm.config[i+1:]...)
			break
		}
	}
	// 从hosts列表中删除
	for i, hostsData := range cm.hostsList {
		if hostsData.ID == id {
			// 删除数据文件
			dataFile := filepath.Join(cm.basePath, "data", "collection", "hosts", "data", hostsData.IDNumber+".json")
			if err := os.Remove(dataFile); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("删除数据文件失败: %w", err)
			}
			cm.hostsList = append(cm.hostsList[:i], cm.hostsList[i+1:]...)
			break
		}
	}
	// 从ID列表中删除
	for i, idNumber := range cm.ids {
		// 找到对应的ID编号
		dataFile := filepath.Join(cm.basePath, "data", "collection", "hosts", "data", idNumber+".json")
		data, err := os.ReadFile(dataFile)
		if err != nil {
			continue
		}
		var hostsData HostsData
		if json.Unmarshal(data, &hostsData) == nil && hostsData.ID == id {
			cm.ids = append(cm.ids[:i], cm.ids[i+1:]...)
			break
		}
	}
	// 保存配置文件
	if err := cm.saveConfigToFile(); err != nil {
		return err
	}
	// 保存ID列表文件
	if err := cm.saveIdsToFile(); err != nil {
		return err
	}
	return nil
}

func (cm *ConfigManager) loadJSONFile(o any, element ...string) error {
	filename := filepath.Join(element...)
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}
