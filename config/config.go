package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	OpenAIModel   string `envconfig:"OPENAI_MODEL" default:"gpt-4o-2024-08-06"`
	AudioLanguage string `envconfig:"AUDIO_LANGUAGE" default:"en"`
	WithImage     bool   `envconfig:"WITH_IMAGE" default:"false"`
	TempDir       string `envconfig:"TEMPDIR" default:"auto"`
}

var ConfigInstance *Config

func init() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}
	if cfg.TempDir == "auto" {
		cfg.TempDir, err = os.MkdirTemp("", "gptslideshow-*")
		if err != nil {
			log.Fatalf("Failed to process environment variables: %v", err)
		}
	}
	ConfigInstance = &cfg
}

// Help prints the environment variables and their default values to stdout.
func Help() {
	err := envconfig.Usage("", &Config{})
	if err != nil {
		log.Printf("Failed to generate help: %v\n", err)
		os.Exit(1)
	}
}
