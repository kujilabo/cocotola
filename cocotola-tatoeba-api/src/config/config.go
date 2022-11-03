package config

import (
	"os"

	"gopkg.in/yaml.v2"

	libconfig "github.com/kujilabo/cocotola/lib/config"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type AppConfig struct {
	Name        string `yaml:"name" validate:"required"`
	HTTPPort    int    `yaml:"httpPort" validate:"required"`
	MetricsPort int    `yaml:"metricsPort" validate:"required"`
}

type AuthConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type ShutdownConfig struct {
	TimeSec1 int `yaml:"timeSec1" validate:"gte=1"`
	TimeSec2 int `yaml:"timeSec2" validate:"gte=1"`
}

type DebugConfig struct {
	GinMode bool `yaml:"ginMode"`
	Wait    bool `yaml:"wait"`
}

type Config struct {
	App      *AppConfig               `yaml:"app" validate:"required"`
	DB       *libconfig.DBConfig      `yaml:"db" validate:"required"`
	Auth     *AuthConfig              `yaml:"auth" validate:"required"`
	Trace    *libconfig.TraceConfig   `yaml:"trace" validate:"required"`
	CORS     *libconfig.CORSConfig    `yaml:"cors" validate:"required"`
	Shutdown *ShutdownConfig          `yaml:"shutdown" validate:"required"`
	Log      *libconfig.LogConfig     `yaml:"log" validate:"required"`
	Swagger  *libconfig.SwaggerConfig `yaml:"swagger" validate:"required"`
	Debug    *DebugConfig             `yaml:"debug"`
}

func LoadConfig(env string) (*Config, error) {
	confContent, err := os.ReadFile("./configs/" + env + ".yml")
	if err != nil {
		return nil, err
	}

	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		return nil, err
	}

	if err := libD.Validator.Struct(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
