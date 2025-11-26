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

	checkBoxList         []*vcl.TCheckBox
	MemoHosts            *vcl.TMemo
	ScrollBox            *vcl.TScrollBox
	StatusBar            *vcl.TStatusBar
	ButtonSystemHosts    *vcl.TButton
	ButtonAddConfigEntry *vcl.TButton
	ButtonShowInfo       *vcl.TButton

	currentEditID string
}

func (f *TFormMain) OnFormCreate(sender vcl.IObject) {
	f.initComponents()
	f.refreshUI()
}

func (f *TFormMain) initComponents() {
	f.SetCaption("GoSwitchHosts v1.0")
	f.SetWidth(1200)
	f.SetHeight(600)
	f.SetPosition(types.PoScreenCenter)
	f.SetDoubleBuffered(true)
	mainPanel := vcl.NewPanel(f)
	mainPanel.SetParent(f)
	mainPanel.SetAlign(types.AlClient)
	mainPanel.SetBevelOuter(types.BvNone)
	leftPanel := vcl.NewPanel(mainPanel)
	leftPanel.SetParent(mainPanel)
	leftPanel.SetWidth(280)
	leftPanel.SetAlign(types.AlLeft)
	leftPanel.SetBevelOuter(types.BvNone)
	topPanel := vcl.NewPanel(leftPanel)
	topPanel.SetParent(leftPanel)
	topPanel.SetHeight(42)
	topPanel.SetAlign(types.AlTop)
	topPanel.SetBevelOuter(types.BvNone)
	f.ButtonAddConfigEntry = vcl.NewButton(topPanel)
	f.ButtonAddConfigEntry.SetParent(topPanel)
	f.ButtonAddConfigEntry.SetCaption("新增配置")
	f.ButtonAddConfigEntry.SetLeft(10)
	f.ButtonAddConfigEntry.SetTop(6)
	f.ButtonAddConfigEntry.SetWidth(80)
	f.ButtonAddConfigEntry.SetHeight(28)
	f.ButtonAddConfigEntry.SetOnClick(func(sender vcl.IObject) { f.onButtonAddConfigClick() })
	f.ScrollBox = vcl.NewScrollBox(leftPanel)
	f.ScrollBox.SetParent(leftPanel)
	f.ScrollBox.SetAlign(types.AlClient)
	rightPanel := vcl.NewPanel(mainPanel)
	rightPanel.SetParent(mainPanel)
	rightPanel.SetAlign(types.AlClient)
	rightPanel.SetBevelOuter(types.BvNone)
	rightTopPanel := vcl.NewPanel(rightPanel)
	rightTopPanel.SetParent(rightPanel)
	rightTopPanel.SetHeight(42)
	rightTopPanel.SetAlign(types.AlTop)
	rightTopPanel.SetBevelOuter(types.BvNone)
	f.ButtonSystemHosts = vcl.NewButton(rightTopPanel)
	f.ButtonSystemHosts.SetParent(rightTopPanel)
	f.ButtonSystemHosts.SetCaption("系统 Hosts")
	f.ButtonSystemHosts.SetLeft(0)
	f.ButtonSystemHosts.SetTop(6)
	f.ButtonSystemHosts.SetWidth(80)
	f.ButtonSystemHosts.SetHeight(28)
	f.ButtonSystemHosts.SetOnClick(func(sender vcl.IObject) { f.onButtonSystemHostsClick() })
	f.ButtonShowInfo = vcl.NewButton(rightTopPanel)
	f.ButtonShowInfo.SetParent(rightTopPanel)
	f.ButtonShowInfo.SetCaption("查看信息")
	f.ButtonShowInfo.SetLeft(96)
	f.ButtonShowInfo.SetTop(6)
	f.ButtonShowInfo.SetWidth(80)
	f.ButtonShowInfo.SetHeight(28)
	f.ButtonShowInfo.SetOnClick(func(sender vcl.IObject) { f.onButtonShowInfoClick() })
	f.MemoHosts = vcl.NewMemo(rightPanel)
	f.MemoHosts.SetParent(rightPanel)
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
	f.updateStatusBar("就绪")
}

func (f *TFormMain) onButtonSystemHostsClick() {
	f.currentEditID = ""
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
	if f.currentEditID == "" {
		f.onButtonSystemHostsClick()
	}
	f.updateStatusBar(fmt.Sprintf("已更新配置: %s", checkBox.Caption()))
}

func (f *TFormMain) onEditLabelClick(config ConfigEntry) {
	f.currentEditID = config.ID
	content, found := configManager.GetHostsByID(config.ID)
	if !found {
		vcl.ShowMessage(fmt.Sprintf("未找到配置项: %s", config.Title))
		return
	}
	f.MemoHosts.SetReadOnly(false)
	f.MemoHosts.SetText(content)
	f.updateStatusBar(fmt.Sprintf("正在编辑: %s", config.Title))
}

func (f *TFormMain) onDeleteLabelClick(config ConfigEntry) {
	if vcl.MessageDlg(fmt.Sprintf("确定要删除配置项 '%s' 吗？", config.Title), types.MtConfirmation, types.MbYes, types.MbNo) == types.MrYes {
		if err := configManager.DeleteConfig(config.ID); err != nil {
			vcl.ShowMessage(fmt.Sprintf("删除配置失败: %v", err))
			return
		}
		if err := configManager.SaveToFile(systemHostsFilename); err != nil {
			vcl.ShowMessage(fmt.Sprintf("保存hosts文件失败: %v", err))
			return
		}
		// 如果正在编辑被删除的配置，清空编辑器
		if f.currentEditID == config.ID {
			f.currentEditID = ""
			f.onButtonSystemHostsClick()
		}
		// 重新加载UI
		f.refreshUI()
		f.updateStatusBar(fmt.Sprintf("已删除配置: %s", config.Title))
	}
}

func (f *TFormMain) onButtonAddConfigClick() {
	title := vcl.InputBox("新增配置", "请输入配置名称:", "")
	if title == "" {
		return
	}
	// 创建新配置
	newConfig := ConfigEntry{
		ID:    uuid.New().String(),
		Title: title,
		On:    false,
	}
	// 保存新配置，内容为空
	if err := configManager.AddConfig(newConfig, ""); err != nil {
		vcl.ShowMessage(fmt.Sprintf("添加配置失败: %v", err))
		return
	}
	if err := configManager.SaveToFile(systemHostsFilename); err != nil {
		vcl.ShowMessage(fmt.Sprintf("保存hosts文件失败: %v", err))
		return
	}
	// 重新加载UI
	f.refreshUI()
	f.updateStatusBar(fmt.Sprintf("已新增配置: %s", title))
}

func (f *TFormMain) onMemoKeyDown(key *types.Char, shift types.TShiftState) {
	if f.currentEditID == "" {
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
	if f.currentEditID == "" {
		return
	}
	content := f.MemoHosts.Text()
	if err := configManager.UpdateHostsContent(f.currentEditID, content); err != nil {
		vcl.ShowMessage(fmt.Sprintf("保存内容失败: %v", err))
		return
	}
	if err := configManager.SaveToFile(systemHostsFilename); err != nil {
		vcl.ShowMessage(fmt.Sprintf("生成hosts文件失败: %v", err))
		return
	}
	f.updateStatusBar(fmt.Sprintf("保存文件: %s", f.currentEditID))
}

func (f *TFormMain) toggleCommentCurrentLine() {
	if f.currentEditID == "" {
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

func (f *TFormMain) refreshUI() {
	// 保存当前编辑状态
	editingID := f.currentEditID
	// 释放并重新创建ScrollBox
	if f.ScrollBox != nil {
		parent := f.ScrollBox.Parent()
		if parent != nil {
			f.ScrollBox = vcl.NewScrollBox(parent)
			f.ScrollBox.SetParent(parent)
			f.ScrollBox.SetAlign(types.AlClient)
		}
	}
	// 重新加载配置并完全重建UI
	configs := configManager.GetConfig()
	var (
		checkBoxTop int32 = 6
		labelTop    int32 = 8
	)
	f.checkBoxList = make([]*vcl.TCheckBox, 0, len(configs))
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
		f.checkBoxList = append(f.checkBoxList, checkBox)
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
		f.currentEditID = ""
		f.onButtonSystemHostsClick()
	}
}

func (f *TFormMain) updateStatusBar(text string) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	f.StatusBar.SetSimpleText(fmt.Sprintf("[%s] %s", now, text))
}
