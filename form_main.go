package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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
	f.loadConfigToUI()
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
	f.ButtonAddConfigEntry.SetOnClick(func(sender vcl.IObject) { f.onButtonShowInfoClick() })
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

func (f *TFormMain) loadConfigToUI() {
	var checkBoxTop int32 = 10
	configs := configManager.GetConfig()
	for i, config := range configs {
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
		editLabel := vcl.NewLabel(f.ScrollBox)
		editLabel.SetParent(f.ScrollBox)
		editLabel.SetCaption("编辑")
		editLabel.SetTop(checkBoxTop + int32(i*34))
		editLabel.SetLeft(212)
		editLabel.SetWidth(40)
		editLabel.SetHeight(32)
		editLabel.Font().SetColor(0xE16941)
		editLabel.SetCursor(types.CrHandPoint)
		editLabel.SetOnClick(func(sender vcl.IObject) { f.onEditLabelClick(config) })
		deleteLabel := vcl.NewLabel(f.ScrollBox)
		deleteLabel.SetParent(f.ScrollBox)
		deleteLabel.SetCaption("删除")
		deleteLabel.SetTop(checkBoxTop + int32(i*34))
		deleteLabel.SetLeft(245)
		deleteLabel.SetWidth(40)
		deleteLabel.SetHeight(32)
		deleteLabel.Font().SetColor(colors.ClRed)
		deleteLabel.SetCursor(types.CrHandPoint)
		deleteLabel.SetOnClick(func(sender vcl.IObject) { f.onDeleteLabelClick(config) })
	}
	f.onButtonSystemHostsClick()
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
	// f.currentEditID = config.ID
	// content, found := configManager.GetHostsByID(config.ID)
	// if !found {
	// 	vcl.ShowMessage(fmt.Sprintf("未找到配置项: %s", config.Title))
	// 	return
	// }
	// f.MemoHosts.SetReadOnly(false)
	// f.MemoHosts.SetText(content)
	// f.updateStatusBar(fmt.Sprintf("正在编辑: %s", config.Title))
}

func (f *TFormMain) onMemoKeyDown(key *types.Char, shift types.TShiftState) {
	if f.currentEditID == "" {
		return
	}
	if shift.In(types.SsCtrl) {
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

func (f *TFormMain) updateStatusBar(text string) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	f.StatusBar.SetSimpleText(fmt.Sprintf("[%s] %s", now, text))
}
