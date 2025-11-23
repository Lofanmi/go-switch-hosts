# CLAUDE.md

## 前言

- 请你使用简体中文，不要使用其它语言
- 在我未同意之前, 你不能提交代码，不能推送到远程仓库
- 不需要生成单元测试、实例、README.md，除非我要求你生成
- 我要求生成单元测试，一般不需要生成示例和文档，除非我要求你生成
- 项目存在 Makefile，编译项目使用 `make build`，而不是 `go build`；依赖更新使用 `make wire`。
- 函数内部禁止空行，保证代码紧凑，必要可加注释，但简单逻辑禁止注释。

## 特殊的代码风格和提交规范

#### 注释需要包含标识符，和Go标准库的文档一样

如 `// String xxx...` 对应方法 `String() string`。

#### 函数返回时，尽量减少显式的结构体初始化，更加简洁，并注意错误处理。

```go
type Data struct{
	A string
}
func A() (res Data, err error) {
	return // 这样写比较简洁
	// 不要这样写
	// return Data{}, errors.New("error")
	// return Data{}, nil
}
```

#### git commit 格式： feat|fix|docs(xxx): message

1. feat(功能点或者是模块点，英文): 功能点是什么
2. fix(修复点或者是模块点，英文): 修复了什么
3. feat(refactor): 重构了什么
注：feat|fix|docs，必须是这三者之一，必须带括号说明改动点。

## 项目信息

### 项目概述
GoSwitchHosts 是基于 Go + govcl 框架开发的轻量级跨平台 hosts 管理工具，旨在提供高性能的原生桌面应用体验，作为 Electron 版本 SwitchHosts 的精简替代方案。

### 核心设计理念
- **性能优先**: 7MB 体积，60MB 内存，秒级启动
- **功能精简**: 仅支持本地 hosts 管理，不做复杂功能
- **原生体验**: 使用 govcl 原生 UI 框架，响应迅速
- **自包含**: 内置 liblcl 运行时，无外部依赖

### 技术架构
- **GUI**: govcl (Go VCL bindings) - 跨平台原生桌面应用
- **存储**: JSON 配置文件存储于 `~/.SwitchHosts` 目录
- **平台**: 通过 build tags 实现 Windows/Linux/macOS 适配
- **构建**: Makefile 支持平台自动检测和多目标编译

### 构建命令
- `make build` - 编译当前平台版本（自动检测）
- `make linux/darwin/windows` - 编译指定平台版本
- `make all` - 编译所有平台版本
- `make clean` - 清理编译产物

### 环境变量
- `GOSH_SWITCHHOSTSDIR` - 自定义配置目录
- `GOSH_HOSTSFILENAME` - 自定义 hosts 文件路径

### 项目限制
- 不支持远程 hosts 配置
- 不支持多文件夹管理
- 不支持代理设置
- 不支持语法高亮（macOS 限制）
- 不支持历史备份