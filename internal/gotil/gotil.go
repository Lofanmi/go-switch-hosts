package gotil

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	EnvGoSwitchHostsLogLevel   = "GO_SWITCH_HOSTS_LOG_LEVEL"
	EnvGoSwitchHostsConfigPath = "GO_SWITCH_HOSTS_CONFIG_PATH"
	EnvGoSwitchHostsConfigName = "GO_SWITCH_HOSTS_CONFIG_NAME"
	EnvGoSwitchHostsConfigType = "GO_SWITCH_HOSTS_CONFIG_TYPE"
)

func StringCut(s, begin, end string, withBegin bool) string {
	beginPos := strings.Index(s, begin)
	if beginPos == -1 {
		return ""
	}
	s = s[beginPos+len(begin):]
	endPos := strings.Index(s, end)
	if endPos == -1 {
		return ""
	}
	result := s[:endPos]
	if withBegin {
		return begin + result
	} else {
		return result
	}
}

func IsDevelopment() bool {
	return log.GetLevel() == log.DebugLevel
}

func Env(key, defaultValue string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = defaultValue
	}
	log.WithField("key", key).WithField("value", value).Debug("读取环境变量")
	return value
}
