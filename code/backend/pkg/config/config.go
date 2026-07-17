package config

import (
	"fmt"
	"strings"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis RedisConfig `mapstructure:"redis"`
	JWT JWTConfig `mapstructure:"jwt"`
	Security SecurityConfig `mapstructure:"security"`
	PII PIIConfig `mapstructure:"pii"`
	CORS CORSConfig `mapstructure:"cors"`
	Log LogConfig `mapstructure:"log"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	User string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode string `mapstructure:"sslmode"`
	MaxConns int32 `mapstructure:"max_conns"`
	MinConns int32 `mapstructure:"min_conns"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB int `mapstructure:"db"`
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	AccessTTL int `mapstructure:"access_ttl"`
	RefreshTTL int `mapstructure:"refresh_ttl"`
	Issuer string `mapstructure:"issuer"`
}

type SecurityConfig struct {
	BcryptCost int `mapstructure:"bcrypt_cost"`
	VerificationCodeTTL int `mapstructure:"verification_code_ttl"`
	VerificationCodeLength int `mapstructure:"verification_code_length"`
	RateLimitThreshold int `mapstructure:"rate_limit_threshold"`
	RateLimitWindow int `mapstructure:"rate_limit_window"`
	RateLimitLock int `mapstructure:"rate_limit_lock"`
	ResetRateLimit int `mapstructure:"reset_rate_limit"`
	DeactivateGraceDays int `mapstructure:"deactivate_grace_days"`
}

type PIIConfig struct {
	EncryptionKey string `mapstructure:"encryption_key"`
	IndexKey string `mapstructure:"index_key"`
}

type CORSConfig struct {
	Origins []string `mapstructure:"origins"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

var cfg *Config

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()
	v.SetEnvPrefix("XY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	cfg = &c
	return &c, nil
}

func Get() *Config {
	if cfg == nil {
		panic("config not loaded")
	}
	return cfg
}
