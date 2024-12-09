package configs

import (
	_ "embed"
	"log"

	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/config"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var cfg []byte

func MustLoadConfig() *config.Config {
	appCfg := &config.Config{}

	if err := yaml.Unmarshal(cfg, &appCfg); err != nil {
		log.Fatalln(err)
	}

	return appCfg
}
