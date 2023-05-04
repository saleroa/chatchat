package config

type Server struct {
	Mode string `mapstructure:"mode" yaml:"mode"`
	Port string `mapstructure:"port" yaml:"port"`
}

func (s Server) Addr() string {
	return ":" + s.Port
}
