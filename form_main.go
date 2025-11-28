package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
	"github.com/ying32/govcl/vcl/types/keys"
)

type TFormMain struct {
	*vcl.TForm
	CurrentEditID        string
	PendingAdd           ConfigEntry
	PendingDelete        ConfigEntry
	CheckBoxList         []*vcl.TCheckBox
	MemoHosts            *vcl.TMemo
	ScrollBox            *vcl.TScrollBox
	StatusBar            *vcl.TStatusBar
	ButtonSystemHosts    *vcl.TButton
	ButtonAddConfigEntry *vcl.TButton
	ButtonShowInfo       *vcl.TButton
	TopPanel             *vcl.TPanel
	MainPanel            *vcl.TPanel
	LeftPanel            *vcl.TPanel
	RightPanel           *vcl.TPanel
	RightTopPanel        *vcl.TPanel
	AddTimer             *vcl.TTimer
	DeleteTimer          *vcl.TTimer
}

func (f *TFormMain) OnFormCreate(sender vcl.IObject) {
	f.initComponents()
	f.refreshConfigEntryUI()
}

func (f *TFormMain) initComponents() {
	f.SetCaption("GoSwitchHosts v1.0")
	f.SetWidth(1200)
	f.SetHeight(600)
	f.SetPosition(types.PoScreenCenter)
	f.SetDoubleBuffered(true)
	f.MainPanel = vcl.NewPanel(f)
	f.MainPanel.SetParent(f)
	f.MainPanel.SetAlign(types.AlClient)
	f.MainPanel.SetBevelOuter(types.BvNone)
	f.LeftPanel = vcl.NewPanel(f.MainPanel)
	f.LeftPanel.SetParent(f.MainPanel)
	f.LeftPanel.SetWidth(280)
	f.LeftPanel.SetAlign(types.AlLeft)
	f.LeftPanel.SetBevelOuter(types.BvNone)
	f.TopPanel = vcl.NewPanel(f.LeftPanel)
	f.TopPanel.SetParent(f.LeftPanel)
	f.TopPanel.SetHeight(42)
	f.TopPanel.SetAlign(types.AlTop)
	f.TopPanel.SetBevelOuter(types.BvNone)
	f.ButtonAddConfigEntry = vcl.NewButton(f.TopPanel)
	f.ButtonAddConfigEntry.SetParent(f.TopPanel)
	f.ButtonAddConfigEntry.SetCaption("新增配置")
	f.ButtonAddConfigEntry.SetLeft(10)
	f.ButtonAddConfigEntry.SetTop(6)
	f.ButtonAddConfigEntry.SetWidth(80)
	f.ButtonAddConfigEntry.SetHeight(28)
	f.ButtonAddConfigEntry.SetOnClick(func(sender vcl.IObject) { f.onButtonAddConfigClick() })
	f.RightPanel = vcl.NewPanel(f.MainPanel)
	f.RightPanel.SetParent(f.MainPanel)
	f.RightPanel.SetAlign(types.AlClient)
	f.RightPanel.SetBevelOuter(types.BvNone)
	f.RightTopPanel = vcl.NewPanel(f.RightPanel)
	f.RightTopPanel.SetParent(f.RightPanel)
	f.RightTopPanel.SetHeight(42)
	f.RightTopPanel.SetAlign(types.AlTop)
	f.RightTopPanel.SetBevelOuter(types.BvNone)
	f.ButtonSystemHosts = vcl.NewButton(f.RightTopPanel)
	f.ButtonSystemHosts.SetParent(f.RightTopPanel)
	f.ButtonSystemHosts.SetCaption("系统 Hosts")
	f.ButtonSystemHosts.SetLeft(0)
	f.ButtonSystemHosts.SetTop(6)
	f.ButtonSystemHosts.SetWidth(85)
	f.ButtonSystemHosts.SetHeight(28)
	f.ButtonSystemHosts.SetOnClick(func(sender vcl.IObject) { f.onButtonSystemHostsClick() })
	f.ButtonShowInfo = vcl.NewButton(f.RightTopPanel)
	f.ButtonShowInfo.SetParent(f.RightTopPanel)
	f.ButtonShowInfo.SetCaption("查看信息")
	f.ButtonShowInfo.SetLeft(96)
	f.ButtonShowInfo.SetTop(6)
	f.ButtonShowInfo.SetWidth(85)
	f.ButtonShowInfo.SetHeight(28)
	f.ButtonShowInfo.SetOnClick(func(sender vcl.IObject) { f.onButtonShowInfoClick() })
	f.MemoHosts = vcl.NewMemo(f.RightPanel)
	f.MemoHosts.SetParent(f.RightPanel)
	f.MemoHosts.SetAlign(types.AlClient)
	f.MemoHosts.SetScrollBars(types.SsBoth)
	f.MemoHosts.SetWordWrap(false)
	f.MemoHosts.SetReadOnly(true)
	f.MemoHosts.SetParentFont(false)
	f.MemoHosts.Font().SetName("Consolas")
	f.MemoHosts.Font().SetSize(10)
	f.MemoHosts.SetOnKeyDown(func(sender vcl.IObject, key *types.Char, shift types.TShiftState) { f.onMemoKeyDown(key, shift) })
	f.StatusBar = vcl.NewStatusBar(f)
	f.StatusBar.SetParent(f)
	f.StatusBar.SetAlign(types.AlBottom)
	f.StatusBar.SetSimplePanel(true)
	f.AddTimer = vcl.NewTimer(f)
	f.AddTimer.SetInterval(50)
	f.AddTimer.SetEnabled(false)
	f.AddTimer.SetOnTimer(func(sender vcl.IObject) { f.onAddTimer() })
	f.DeleteTimer = vcl.NewTimer(f)
	f.DeleteTimer.SetInterval(50)
	f.DeleteTimer.SetEnabled(false)
	f.DeleteTimer.SetOnTimer(func(sender vcl.IObject) { f.onDeleteTimer() })
	f.updateStatusBar("就绪")
	f.AutoAdjustLayout(types.LapAutoAdjustForDPI, f.DesignTimePPI(), vcl.Screen.PixelsPerInch(), 0, 0)
}

func (f *TFormMain) onButtonSystemHostsClick() {
	f.CurrentEditID = ""
	data, err := os.ReadFile(systemHostsFilename)
	if err != nil {
		f.updateStatusBar(fmt.Sprintf("# 无法读取系统 hosts 文件: %v, 错误: %v", systemHostsFilename, err))
	} else {
		f.MemoHosts.SetText(string(data))
		f.MemoHosts.SetReadOnly(true)
	}
}

func (f *TFormMain) onButtonShowInfoClick() {
	vcl.MessageDlg(fmt.Sprintf("配置路径: %s%s%s系统 hosts 文件路径: %s%s%s%s作者: %s%s%sGitHub: %s%s%s",
		LineEnding, configSwitchHostsDir, LineEnding,
		LineEnding, systemHostsFilename, LineEnding,
		LineEnding,
		LineEnding, "Lofanmi", LineEnding,
		LineEnding, "https://github.com/Lofanmi/go-switch-hosts", LineEnding,
	), types.MtInformation)
}

func (f *TFormMain) onCheckBoxClick(checkBox *vcl.TCheckBox, configID string) {
	if err := configManager.UpdateConfigStatus(configID, checkBox.Checked()); err != nil {
		vcl.ShowMessage(fmt.Sprintf("更新配置状态失败: %v", err))
		return
	}
	if err := configManager.SaveToFile(systemHostsFilename); err != nil {
		vcl.ShowMessage(fmt.Sprintf("保存文件失败: %v", err))
		return
	}
	if f.CurrentEditID == "" {
		f.onButtonSystemHostsClick()
	}
	f.updateStatusBar(fmt.Sprintf("已更新配置: %s", checkBox.Caption()))
}

func (f *TFormMain) onEditLabelClick(config ConfigEntry) {
	f.CurrentEditID = config.ID
	content, found := configManager.GetHostsByID(config.ID)
	if !found {
		vcl.ShowMessage(fmt.Sprintf("未找到配置项: %s", config.Title))
		return
	}
	f.MemoHosts.SetReadOnly(false)
	f.MemoHosts.SetText(content)
	f.updateStatusBar(fmt.Sprintf("正在编辑: %s", config.Title))
}

func (f *TFormMain) onButtonAddConfigClick() {
	title := vcl.InputBox("新增配置", "请输入配置名称:", "")
	if title == "" {
		return
	}
	f.PendingAdd = ConfigEntry{
		ID:    uuid.New().String(),
		Title: title,
		On:    false,
	}
	f.AddTimer.SetEnabled(true)
}

func (f *TFormMain) onDeleteLabelClick(config ConfigEntry) {
	if vcl.MessageDlg(fmt.Sprintf("确定要删除配置项 '%s' 吗？", config.Title), types.MtConfirmation, types.MbYes, types.MbNo) == types.MrYes {
		f.PendingDelete = config
		f.DeleteTimer.SetEnabled(true)
	}
}

func (f *TFormMain) onMemoKeyDown(key *types.Char, shift types.TShiftState) {
	if f.CurrentEditID == "" {
		return
	}
	if shift.In(CommandKeyCode) {
		hasSelection := f.MemoHosts.SelText() != ""
		switch *key {
		case keys.VkS:
			f.saveCurrentContent()
		case keys.VkSlash, keys.VkDivide:
			if !hasSelection {
				f.toggleCommentCurrentLine()
				f.saveCurrentContent()
			}
		case keys.VkC:
			if !hasSelection {
				f.copyCurrentLine()
			}
		case keys.VkX:
			if !hasSelection {
				f.cutCurrentLine()
			}
		case keys.VkV:
			if !hasSelection {
				f.pasteAtCurrentLine()
			}
		}
	}
}

func (f *TFormMain) saveCurrentContent() {
	if f.CurrentEditID == "" {
		return
	}
	content := f.MemoHosts.Text()
	if err := configManager.UpdateHostsContent(f.CurrentEditID, content); err != nil {
		vcl.ShowMessage(fmt.Sprintf("保存内容失败: %v", err))
		return
	}
	if err := configManager.SaveToFile(systemHostsFilename); err != nil {
		vcl.ShowMessage(fmt.Sprintf("生成hosts文件失败: %v", err))
		return
	}
	f.updateStatusBar(fmt.Sprintf("保存文件: %s", f.CurrentEditID))
}

func (f *TFormMain) toggleCommentCurrentLine() {
	if f.CurrentEditID == "" {
		return
	}
	lineNum := f.MemoHosts.CaretPos().Y
	lineText := strings.TrimSpace(f.MemoHosts.Lines().S(lineNum))
	if lineText == "" {
		return
	}
	f.MemoHosts.Lines().BeginUpdate()
	if len(lineText) > 0 && (lineText[0] == '#') {
		f.MemoHosts.Lines().SetS(lineNum, strings.TrimSpace(lineText[1:]))
	} else {
		f.MemoHosts.Lines().SetS(lineNum, "# "+strings.TrimSpace(lineText))
	}
	f.MemoHosts.Lines().EndUpdate()
}

func (f *TFormMain) copyCurrentLine() {
	lineNum := f.MemoHosts.CaretPos().Y
	lineText := f.MemoHosts.Lines().S(lineNum)
	vcl.Clipboard.SetAsText(LineEnding + lineText)
	f.updateStatusBar("已复制当前行")
}

func (f *TFormMain) cutCurrentLine() {
	lineNum := f.MemoHosts.CaretPos().Y
	lineText := f.MemoHosts.Lines().S(lineNum)
	vcl.Clipboard.SetAsText(LineEnding + lineText)
	f.MemoHosts.Lines().BeginUpdate()
	f.MemoHosts.Lines().Delete(lineNum)
	f.MemoHosts.Lines().EndUpdate()
	f.saveCurrentContent()
	f.updateStatusBar("已剪切当前行")
}

func (f *TFormMain) pasteAtCurrentLine() {
	f.saveCurrentContent()
	f.updateStatusBar("已粘贴到当前行")
}

func (f *TFormMain) refreshConfigEntryUI() {
	// 保存当前编辑状态
	editingID := f.CurrentEditID
	// 释放并重新创建ScrollBox
	if f.ScrollBox != nil {
		f.ScrollBox.Free()
	}
	f.ScrollBox = vcl.NewScrollBox(nil)
	f.ScrollBox.SetParent(f.LeftPanel)
	f.ScrollBox.SetAlign(types.AlClient)
	defer f.ScrollBox.AutoAdjustLayout(types.LapAutoAdjustForDPI, f.DesignTimePPI(), vcl.Screen.PixelsPerInch(), 0, 0)
	// 重新加载配置并完全重建UI
	configs := configManager.GetConfig()
	var (
		checkBoxTop int32 = 6
		labelTop    int32 = 8
	)
	f.CheckBoxList = make([]*vcl.TCheckBox, 0, len(configs))
	for i, config := range configs {
		// 创建CheckBox
		checkBox := vcl.NewCheckBox(f.ScrollBox)
		checkBox.SetParent(f.ScrollBox)
		checkBox.SetCaption(config.Title)
		checkBox.SetChecked(config.On)
		checkBox.SetTop(checkBoxTop + int32(i*34))
		checkBox.SetLeft(10)
		checkBox.SetWidth(220)
		checkBox.SetHeight(32)
		checkBox.SetOnClick(func(sender vcl.IObject) { f.onCheckBoxClick(checkBox, config.ID) })
		f.CheckBoxList = append(f.CheckBoxList, checkBox)
		// 创建编辑Label
		editLabel := vcl.NewLabel(f.ScrollBox)
		editLabel.SetParent(f.ScrollBox)
		editLabel.SetCaption("编辑")
		editLabel.SetTop(labelTop + int32(i*34))
		editLabel.SetLeft(206)
		editLabel.SetWidth(40)
		editLabel.SetHeight(32)
		editLabel.Font().SetColor(0xE16941)
		editLabel.SetCursor(types.CrHandPoint)
		editLabel.SetOnClick(func(sender vcl.IObject) { f.onEditLabelClick(config) })
		// 创建删除Label
		deleteLabel := vcl.NewLabel(f.ScrollBox)
		deleteLabel.SetParent(f.ScrollBox)
		deleteLabel.SetCaption("删除")
		deleteLabel.SetTop(labelTop + int32(i*34))
		deleteLabel.SetLeft(240)
		deleteLabel.SetWidth(40)
		deleteLabel.SetHeight(32)
		deleteLabel.Font().SetColor(colors.ClRed)
		deleteLabel.SetCursor(types.CrHandPoint)
		deleteLabel.SetOnClick(func(sender vcl.IObject) { f.onDeleteLabelClick(config) })
	}
	// 如果没有正在编辑的配置，则显示系统hosts
	if editingID == "" {
		f.onButtonSystemHostsClick()
		return
	}
	// 如果之前正在编辑的配置被删除了，则显示系统hosts
	found := false
	for _, config := range configs {
		if config.ID == editingID {
			found = true
			break
		}
	}
	if !found {
		f.CurrentEditID = ""
		f.onButtonSystemHostsClick()
	}
}

func (f *TFormMain) updateStatusBar(text string) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	f.StatusBar.SetSimpleText(fmt.Sprintf("[%s] %s", now, text))
}

func (f *TFormMain) onAddTimer() {
	var empty ConfigEntry
	f.AddTimer.SetEnabled(false)
	if f.PendingAdd != empty {
		config := f.PendingAdd
		f.PendingAdd = empty
		if err := configManager.AddConfig(config, ""); err != nil {
			vcl.ShowMessage(fmt.Sprintf("添加配置失败: %v", err))
			return
		}
		if err := configManager.SaveToFile(systemHostsFilename); err != nil {
			vcl.ShowMessage(fmt.Sprintf("保存hosts文件失败: %v", err))
			return
		}
		f.refreshConfigEntryUI()
		f.updateStatusBar(fmt.Sprintf("已新增配置: %s", config.Title))
	}
}

func (f *TFormMain) onDeleteTimer() {
	var empty ConfigEntry
	f.DeleteTimer.SetEnabled(false)
	if f.PendingDelete != empty {
		config := f.PendingDelete
		f.PendingDelete = empty
		if err := configManager.DeleteConfig(config.ID); err != nil {
			vcl.ShowMessage(fmt.Sprintf("删除配置失败: %v", err))
			return
		}
		if err := configManager.SaveToFile(systemHostsFilename); err != nil {
			vcl.ShowMessage(fmt.Sprintf("保存hosts文件失败: %v", err))
			return
		}
		if f.CurrentEditID == config.ID {
			f.CurrentEditID = ""
			f.onButtonSystemHostsClick()
		}
		f.refreshConfigEntryUI()
		f.updateStatusBar(fmt.Sprintf("已删除配置: %s", config.Title))
	}
}
