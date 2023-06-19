package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	UserSrvInfo UserSrvConfig `mapstructure:"user-srv"`
}
