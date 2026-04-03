package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	AI       AIConfig       `yaml:"ai"`
	Auth     AuthConfig     `yaml:"auth"`
	Crawl    CrawlConfig    `yaml:"crawl"`
	Cache    CacheConfig    `yaml:"cache"`
}

type ServerConfig struct {
	Port    int    `yaml:"port"`
	DataDir string `yaml:"data_dir"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type AIConfig struct {
	DeepSeekAPIKey string `yaml:"deepseek_api_key"`
	DeepSeekBase   string `yaml:"deepseek_base_url"`
	Model          string `yaml:"model"`
}

type AuthConfig struct {
	JWTSecret       string `yaml:"jwt_secret"`
	JWTExpireHours  int    `yaml:"jwt_expire_hours"`
	SMSProvider     string `yaml:"sms_provider"`      // aliyun or stub
	SMSAccessKey    string `yaml:"sms_access_key"`    // 阿里云 AccessKey ID
	SMSAccessSecret string `yaml:"sms_access_secret"` // 阿里云 AccessKey Secret
	SMSSignName     string `yaml:"sms_sign_name"`     // 短信签名，如"江湖小报"
	SMSTemplateCode string `yaml:"sms_template_code"` // 短信模板ID，如"SMS_123456789"
}

type CrawlConfig struct {
	IntervalMinutes int          `yaml:"interval_minutes"`
	Feeds           []FeedConfig `yaml:"feeds"`
}

type FeedConfig struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Enabled bool   `yaml:"enabled"`
}

type CacheConfig struct {
	MaxCost     int64 `yaml:"max_cost"`
	NumCounters int64 `yaml:"num_counters"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:    8080,
			DataDir: "./data",
		},
		Database: DatabaseConfig{
			Path: "./data/jianghu.db",
		},
		AI: AIConfig{
			DeepSeekBase: "https://api.deepseek.com",
			Model:        "deepseek-chat",
		},
		Auth: AuthConfig{
			JWTExpireHours: 168, // 7 days
		},
		Crawl: CrawlConfig{
			IntervalMinutes: 30,
		},
		Cache: CacheConfig{
			MaxCost:     1 << 28, // 256MB
			NumCounters: 1e7,
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Try environment variables as fallback
		if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
			cfg.AI.DeepSeekAPIKey = key
		}
		if secret := os.Getenv("JWT_SECRET"); secret != "" {
			cfg.Auth.JWTSecret = secret
		}
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Override with environment variables
	if key := os.Getenv("DEEPSEEK_API_KEY"); key != "" {
		cfg.AI.DeepSeekAPIKey = key
	}
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.Auth.JWTSecret = secret
	}

	return cfg, nil
}
