package config

type Logger struct {
	LogLevel   string `mapstructure:"loglevel" yaml:"loglevel"`
	SavePath   string `mapstructure:"savePath" yaml:"savePath"`
	MaxSize    int    `mapstructure:"maxsize" yaml:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups"`
	IsCompress bool   `mapstructure:"isCompress" yaml:"isCompress"`
}
