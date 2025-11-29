# 检测当前操作系统
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# 根据平台设置默认二进制文件名和编译参数
ifeq ($(OS),Windows_NT)
	BINARY = GoSwitchHosts.exe
	LDFLAGS = -ldflags "-s -w -H windowsgui"
	PLATFORM = windows
else ifeq ($(UNAME_S),Linux)
	BINARY = GoSwitchHosts-linux-$(UNAME_M)
	LDFLAGS = -ldflags "-s -w"
	PLATFORM = linux
else ifeq ($(UNAME_S),Darwin)
	BINARY = GoSwitchHosts-darwin-$(UNAME_M)
	LDFLAGS = -ldflags "-s -w"
	PLATFORM = darwin
else
	# 默认为 Windows 平台
	BINARY = GoSwitchHosts.exe
	LDFLAGS = -ldflags "-s -w -H windowsgui"
	PLATFORM = windows
endif

# Windows 资源文件编译
res:
	windres -i GoSwitchHosts.rc -o defaultRes_windows_386.syso   -F pe-i386
	windres -i GoSwitchHosts.rc -o defaultRes_windows_amd64.syso -F pe-x86-64

# 默认构建 - 编译当前平台版本
build: $(PLATFORM)

# Linux 平台编译
linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags tempdll $(LDFLAGS) -o GoSwitchHosts-linux-amd64 .

# macOS 平台编译
darwin:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o GoSwitchHosts-darwin-amd64 .

# macOS APP 打包
macapp: darwin icns
	@echo "开始创建 macOS APP 包..."
	@mkdir -p GoSwitchHosts.app/Contents/MacOS
	@mkdir -p GoSwitchHosts.app/Contents/Resources
	@cp GoSwitchHosts-darwin-amd64 GoSwitchHosts.app/Contents/MacOS/GoSwitchHosts
	@if [ -f macos/liblcl.dylib ]; then \
		cp macos/liblcl.dylib GoSwitchHosts.app/Contents/MacOS/; \
		echo "复制 liblcl.dylib 成功"; \
	else \
		echo "警告: macos/liblcl.dylib 不存在，请确保已复制到macos目录"; \
	fi
	@cp macos/Info.plist GoSwitchHosts.app/Contents/
	@if [ -f macos/GoSwitchHosts.icns ]; then \
		cp macos/GoSwitchHosts.icns GoSwitchHosts.app/Contents/Resources/; \
		echo "复制应用图标成功"; \
	else \
		echo "警告: macos/GoSwitchHosts.icns 不存在"; \
	fi
	@chmod +x GoSwitchHosts.app/Contents/MacOS/GoSwitchHosts
	@echo "GoSwitchHosts.app 包创建完成"
	@echo "应用路径: $(PWD)/GoSwitchHosts.app"

# 生成 icns 图标 (如果不存在)
icns:
	@if [ ! -f macos/GoSwitchHosts.icns ] && [ -f GoSwitchHosts.ico ]; then \
		echo "从 ICO 生成 ICNS 图标..."; \
		mkdir -p GoSwitchHosts.iconset; \
		magick "GoSwitchHosts.ico[1]" GoSwitchHosts.iconset/icon_16x16.png; \
		magick "GoSwitchHosts.ico[2]" GoSwitchHosts.iconset/icon_16x16@2x.png; \
		magick "GoSwitchHosts.ico[4]" GoSwitchHosts.iconset/icon_32x32.png; \
		magick "GoSwitchHosts.ico[5]" GoSwitchHosts.iconset/icon_32x32@2x.png; \
		magick "GoSwitchHosts.ico[0]" GoSwitchHosts.iconset/icon_128x128.png; \
		magick "GoSwitchHosts.ico[3]" GoSwitchHosts.iconset/icon_256x256.png; \
		cp GoSwitchHosts.iconset/icon_256x256.png GoSwitchHosts.iconset/icon_128x128@2x.png; \
		cp GoSwitchHosts.iconset/icon_256x256.png GoSwitchHosts.iconset/icon_256x256@2x.png; \
		cp GoSwitchHosts.iconset/icon_256x256.png GoSwitchHosts.iconset/icon_512x512.png; \
		cp GoSwitchHosts.iconset/icon_256x256.png GoSwitchHosts.iconset/icon_512x512@2x.png; \
		iconutil -c icns GoSwitchHosts.iconset -o macos/GoSwitchHosts.icns; \
		rm -rf GoSwitchHosts.iconset; \
	fi

# Windows 平台编译
windows: res
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -tags tempdll $(LDFLAGS) -o GoSwitchHosts-windows-amd64.exe .

# 编译所有平台版本
all: linux darwin windows macapp

# 清理编译产物
clean:
	rm -f GoSwitchHosts-*
	rm -f GoSwitchHosts.exe
	rm -f *.syso
	rm -f *.icns
	rm -rf GoSwitchHosts.app
	rm -rf GoSwitchHosts.iconset
	rm -f GoSwitchHosts_temp-*.png

.PHONY: res build linux darwin windows macapp icns all clean
