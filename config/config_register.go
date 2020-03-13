package config

type RegisterConfig struct {
	Consul        bool   `kiper_value:"name:consul;help:use consul or not;default:true"`
	ConsulAddress string `kiper_value:"name:consul_address;help:consul server address;default:127.0.0.1"`
	ConsulPort    *Port  `kiper_value:"name:consul_port;help:consul server port;default:8500"`
}

func newRegisterConfig() *RegisterConfig {
	return &RegisterConfig{
		ConsulPort: &Port{},
	}
}
