package ginx

import "github.com/gin-gonic/gin"

const (
	DefaultGinMode            = gin.DebugMode
	DefaultListenAddr         = ":8080"
	DefaultReadTimeout        = 60
	DefaultWriteTimeout       = 60
	DefaultCheckHealthApiPath = "/api/health"
	DefaultLogLevel           = "info"
)

type Config struct {
	GinMode              string // debug,release,test
	Addr                 string
	ReadTimeout          int
	WriteTimeout         int
	EnableCheckHealthApi bool
	CheckHealthApiPath   string
	EnableRequestId      bool
	EncryptSecret        string
	CookieDomain         string
	CookieExpiredSeconds int
}

func loadDefaultConfig(cfg *Config) *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.GinMode == "" {
		cfg.GinMode = DefaultGinMode
	}
	if cfg.Addr == "" {
		cfg.Addr = DefaultListenAddr
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = DefaultReadTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = DefaultWriteTimeout
	}
	if cfg.EnableCheckHealthApi {
		if cfg.CheckHealthApiPath == "" {
			cfg.CheckHealthApiPath = DefaultCheckHealthApiPath
		}
	}
	return cfg
}
