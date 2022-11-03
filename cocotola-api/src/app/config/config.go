package config

import (
	"embed"
	"fmt"
	"os"

	_ "embed"

	"gopkg.in/yaml.v2"

	libconfig "github.com/kujilabo/cocotola/lib/config"
	libD "github.com/kujilabo/cocotola/lib/domain"
)

type AppConfig struct {
	Name          string `yaml:"name" validate:"required"`
	HTTPPort      int    `yaml:"httpPort" validate:"required"`
	MetricsPort   int    `yaml:"metricsPort" validate:"required"`
	OwnerPassword string `yaml:"ownerPassword" validate:"required"`
	TestUserEmail string `yaml:"testUserEmail" validate:"required"`
}

type AuthConfig struct {
	SigningKey          string `yaml:"signingKey"`
	AccessTokenTTLMin   int    `yaml:"accessTokenTtlMin" validate:"gte=1"`
	RefreshTokenTTLHour int    `yaml:"refreshTokenTtlHour" validate:"gte=1"`
	GoogleCallbackURL   string `yaml:"googleCallbackUrl" validate:"required"`
	GoogleClientID      string `yaml:"googleClientId" validate:"required"`
	GoogleClientSecret  string `yaml:"googleClientSecret" validate:"required"`
	APITimeoutSec       int    `yaml:"apiTimeoutSec" validate:"gte=1"`
}

type TranslatorConfig struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	TimeoutSec int    `yaml:"timeoutSec" validate:"gte=1"`
	Username   string `yaml:"username" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
	GRPCAddr   string `yaml:"grpcAddr" validate:"required"`
}

type TatoebaConfig struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	TimeoutSec int    `yaml:"timeoutSec" validate:"gte=1"`
	Username   string `yaml:"username" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
}

type SynthesizerConfig struct {
	Endpoint   string `yaml:"endpoint" validate:"required"`
	TimeoutSec int    `yaml:"timeoutSec" validate:"gte=1"`
	Username   string `yaml:"username" validate:"required"`
	Password   string `yaml:"password" validate:"required"`
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
	App         *AppConfig               `yaml:"app" validate:"required"`
	DB          *libconfig.DBConfig      `yaml:"db" validate:"required"`
	Auth        *AuthConfig              `yaml:"auth" validate:"required"`
	Translator  *TranslatorConfig        `yaml:"translator" validate:"required"`
	Tatoeba     *TatoebaConfig           `yaml:"tatoeba" validate:"required"`
	Synthesizer *SynthesizerConfig       `yaml:"synthesizer" validate:"required"`
	Trace       *libconfig.TraceConfig   `yaml:"trace" validate:"required"`
	CORS        *libconfig.CORSConfig    `yaml:"cors" validate:"required"`
	Shutdown    *ShutdownConfig          `yaml:"shutdown" validate:"required"`
	Log         *libconfig.LogConfig     `yaml:"log" validate:"required"`
	Swagger     *libconfig.SwaggerConfig `yaml:"swagger" validate:"required"`
	Debug       *DebugConfig             `yaml:"debug"`
}

//go:embed local.yml
//go:embed production.yml
var config embed.FS

func LoadConfig(env string) (*Config, error) {
	fmt.Println(os.Getwd())

	confContent, err := config.ReadFile(env + ".yml")
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
