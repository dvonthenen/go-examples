package config

import "flag"

//Config is the representation of the config
type Config struct {
	GatewayIP string
	Username  string
	Password  string
	Version   string
	SdsList   string
}

//AddFlags adds flags to the command line parsing
func (cfg *Config) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&cfg.GatewayIP, "gateway.ip", cfg.GatewayIP, "ScaleIO Gateway IP")
	fs.StringVar(&cfg.Username, "gateway.username", cfg.Username, "ScaleIO Gateway Username")
	fs.StringVar(&cfg.Password, "gateway.password", cfg.Password, "ScaleIO Gateway Password")
	fs.StringVar(&cfg.Version, "gateway.version", cfg.Version, "ScaleIO Gateway Version")
	fs.StringVar(&cfg.SdsList, "gateway.sds", cfg.SdsList, "ScaleIO SDS List")
}

//NewConfig creates a new Config object
func NewConfig() *Config {
	return &Config{
		GatewayIP: env("SCALEIO_GATEWAY_IP", "127.0.0.1"),
		Username:  env("SCALEIO_GATEWAY_USERNAME", "admin"),
		Password:  env("SCALEIO_GATEWAY_PASSWORD", "admin"),
		Version:   env("SCALEIO_GATEWAY_VERSION", ""),
		SdsList:   env("SCALEIO_SDS_LIST", ""),
	}
}
