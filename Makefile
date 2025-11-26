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
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o GoSwitchHosts-linux-amd64 .

# macOS 平台编译
darwin:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -tags tempdll $(LDFLAGS) -o GoSwitchHosts-darwin-amd64 .

# Windows 平台编译
windows: res
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -tags tempdll $(LDFLAGS) -o GoSwitchHosts-windows-amd64.exe .

# 编译所有平台版本
all: linux darwin windows

# 清理编译产物
clean:
	rm -f GoSwitchHosts-*
	rm -f GoSwitchHosts.exe
	rm -f *.syso

.PHONY: res build linux darwin windows all clean