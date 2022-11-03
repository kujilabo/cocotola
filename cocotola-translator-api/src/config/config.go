package config

import (
	"embed"
	"os"

	"gopkg.in/yaml.v2"

	libconfig "github.com/kujilabo/cocotola/lib/config"
	lib "github.com/kujilabo/cocotola/lib/domain"
)

type AppConfig struct {
	Name        string `yaml:"name" validate:"required"`
	HTTPPort    int    `yaml:"httpPort" validate:"required"`
	GRPCPort    int    `yaml:"grpcPort" validate:"required"`
	MetricsPort int    `yaml:"metricsPort" validate:"required"`
}

type HTTPServer struct {
	Port                 int `yaml:"port" validate:"required"`
	ShutdownTimeSec      int `yaml:"shutdownTimeSec" validate:"gte=1"`
	ReadHeaderTimeoutSec int `yaml:"readHeaderTimeoutSec" validate:"gte=1"`
}

type AuthConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type AzureConfig struct {
	SubscriptionKey string `yaml:"subscriptionKey" validate:"required"`
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
	Azure    *AzureConfig             `yaml:"azure" validate:"required"`
	Trace    *libconfig.TraceConfig   `yaml:"trace" validate:"required"`
	CORS     *libconfig.CORSConfig    `yaml:"cors" validate:"required"`
	Shutdown *ShutdownConfig          `yaml:"shutdown" validate:"required"`
	Log      *libconfig.LogConfig     `yaml:"log" validate:"required"`
	Debug    *DebugConfig             `yaml:"debug"`
	Swagger  *libconfig.SwaggerConfig `yaml:"swagger" validate:"required"`
}

//go:embed local.yml
//go:embed production.yml
var config embed.FS

func LoadConfig(env string) (*Config, error) {
	confContent, err := config.ReadFile(env + ".yml")
	if err != nil {
		return nil, err
	}

	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		return nil, err
	}

	if err := lib.Validator.Struct(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
