# GoSwitchHosts

一个基于 Go 语言开发的跨平台 hosts 文件管理工具，提供图形化界面来快速切换不同的 hosts 配置。这是一个轻量级、高性能的原生桌面应用，专为追求极致性能和简洁体验的用户设计。

![GoSwitchHosts](GoSwitchHosts.ico)

## 🌟 项目初衷

本项目可以简要替代 [oldj/SwitchHosts](https://github.com/oldj/SwitchHosts) 项目，解决原版本使用 Electron 框架导致的资源占用过高问题。

但它并不是完全替代，它仅仅支持本地 hosts 的读取、新增和切换，并不会支持远程、文件夹、代理、备份等其它功能，也不会支持语法高亮。

### 🚀 性能对比

| 项目 | 技术栈 | 占用空间     | 运行内存 | 启动速度 | UI 响应 | 劣势 |
|------|--------|----------|---------|----------|---------|------|
| **GoSwitchHosts** | Go + govcl (原生 UI) | **~7MB** | **~60MB** | 极快 | 原生流畅 | 功能精简，仅支持本地 hosts |
| SwitchHosts (原版) | Electron + Web 技术 | ~265MB   | ~630MB | 较慢 | 有卡顿、粘腻感 | 资源占用高，体积大 |

### ✨ 核心优势

- **🎯 轻量级**: 编译后的可执行文件不到 8MB，包含完整的 UI 运行时库
- **⚡ 高性能**: 运行时内存占用仅约 60MB，相比 Electron 版本节省 90% 内存
- **🖥️ 原生体验**: 基于 govcl 框架的原生桌面 UI，响应迅速，无卡顿粘腻感
- **🚀 启动迅速**: 秒级启动，无需等待 Web 环境加载
- **📦 自包含**: 内置 liblcl 运行时，无需额外依赖
- **🔧 精简专注**: 专注于本地 hosts 管理，界面简洁，操作直观

### ⚠️ 功能限制

本项目采用精简设计理念，**并非**完全替代原版 SwitchHosts，目前**不支持**以下功能：

- ❌ 远程 hosts 配置（URL 方式）
- ❌ 多文件夹配置管理
- ❌ 代理设置相关功能
- ❌ hosts 语法高亮显示（govcl 的 TRichEdit 未完全支持 macOS 系统）
- ❌ 历史备份

如需上述高级功能，请继续使用原版 [SwitchHosts](https://github.com/oldj/SwitchHosts)。

## 功能特性

- 🖥️ 图形化界面：基于 govcl 框架的跨平台桌面应用
- 🔄 快速切换：一键切换不同的 hosts 配置
- 📝 配置管理：支持添加、编辑、删除 hosts 配置项
- 🔍 系统集成：直接编辑系统 hosts 文件
- 💾 配置持久化：配置存储在用户目录下，重启不丢失
- 🎯 简洁高效：轻量级设计，运行快速稳定

## 环境要求

- **支持平台**: Windows、Linux、macOS
- Go 1.25+
- **编译依赖**:
  - Windows: GCC (用于编译资源文件)
  - Linux/macOS: GCC (可选，用于交叉编译)

## 编译运行

### 编译项目

```bash
# 编译当前平台版本（自动检测平台）
make build

# 编译特定平台版本
make linux     # 编译 Linux AMD64 版本
make darwin    # 编译 macOS AMD64 版本
make windows   # 编译 Windows AMD64 版本（包含资源文件）

# 编译所有支持的平台
make all

# 清理编译产物
make clean

# 仅编译 Windows 资源文件（如需要）
make res
```

### 环境变量

支持以下环境变量进行自定义配置：

- `GOSH_SWITCHHOSTSDIR`: 自定义配置文件目录（默认：`~/.SwitchHosts`）
- `GOSH_HOSTSFILENAME`: 自定义 hosts 文件路径（默认：系统 hosts 文件路径）

## 项目结构

```
go-switch-hosts/
├── main.go              # 主程序入口
├── form_main.go         # 主窗体界面实现
├── config_manager.go    # 配置管理器
├── utils.go            # 通用工具函数
├── utils_windows.go    # Windows 平台相关函数
├── utils_unix.go       # Unix 平台相关函数
├── GoSwitchHosts.rc    # 资源文件
├── Makefile           # 编译脚本
└── go.mod             # Go 模块文件
```

## 技术栈

- **语言**: Go 1.25+
- **GUI 框架**: [govcl](https://github.com/ying32/govcl) - Go 语言的 VCL 绑定
- **构建工具**: Make + windres

## 致谢

- 图标作者：https://www.iconfont.cn/collections/detail?cid=19977
- 感谢 [govcl](https://github.com/ying32/govcl) 项目提供的优秀 GUI 框架

## 许可证

本项目采用 Apache-2.0 许可证，详见 [LICENSE](LICENSE) 文件。