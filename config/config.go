package config

import (
	"path/filepath"

	"github.com/leeif/kiper"
)

var config *Config

type Config struct {
	Version  string
	Database *DatabaseConfig `kiper_config:"name:database"`
	Log      *LogConfig      `kiper_config:"name:log"`
	Server   *ServerConfig   `kiper_config:"name:server"`
	Register *RegisterConfig `kiper_config:"name:register"`
	Cors     *CorsConfig     `kiper_config:"name:cors"`
}

func NewConfig(args []string, version string) (*Config, error) {
	c := &Config{
		Database: newDatabaseConfig(),
		Log:      newLogConfig(),
		Server:   newServerConfig(),
		Register: newRegisterConfig(),
		Cors:     newCorsConfig(),
	}
	kiper := kiper.NewKiper(filepath.Base(args[0]), "My server")
	kiper.Kingpin.Version(version)
	kiper.Kingpin.HelpFlag.Short('h')

	kiper.SetConfigFileFlag("config.file", "config file", "./config.json")

	if err := kiper.Parse(c, args[1:]); err != nil {
		return nil, err
	}
	c.Version = version
	return c, nil
}
