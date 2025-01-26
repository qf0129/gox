package confx

import (
	"log/slog"
	"os"

	"github.com/qf0129/gox/pkg/jsonx"
)

func LoadJsonConfig(target any) {
	LoadConfig("conf.json", target)
}

func LoadConfig(configFile string, target any) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		slog.Warn("Read config file failed, use default config")
		return
	}
	err = jsonx.Unmarshal(data, target)
	if err != nil {
		panic("UnmarshalConfigFailed: " + err.Error())
	}
	slog.Info("### Load config from " + configFile)
}
