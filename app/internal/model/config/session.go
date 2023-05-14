package config

type Session struct {
	Name      string `mapstructure:"name" yaml:"name"`
	SecretKey string `mapstructure:"secretKey" yaml:"secretKey"`
	MaxAge    int    `mapstructure:"maxAge" yaml:"maxAge"`
}
