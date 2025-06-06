package config

type HTTPServer struct {
	Port int `mapstructure:"port"`
}

type Config struct {
	HTTPServer HTTPServer `mapstructure:"http_server"`
}