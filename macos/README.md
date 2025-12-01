# macOS APP 资源文件

此目录包含构建 macOS 应用包（.app）所需的资源文件：

## 文件说明

- `Info.plist` - macOS 应用配置文件，包含应用元数据和权限配置
- `GoSwitchHosts.icns` - macOS 应用图标文件，包含多分辨率图标
- `liblcl.dylib` - govcl 框架的核心 UI 库文件

## 构建说明

使用 `make app` 命令时会自动从此目录读取这些文件并构建 macOS 应用包。

## 图标生成

如果需要重新生成图标文件，可以运行：
```bash
make icns
```

此命令会从 `GoSwitchHosts.ico` 文件提取原始尺寸图像并生成高质量的 ICNS 文件。