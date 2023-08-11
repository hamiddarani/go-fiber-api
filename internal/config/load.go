package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	delimeter      = "."
	seperator      = "__"
	envPrefix      = "PHONEBOOK_"
	tagName        = "koanf"
	upTemplate     = "================ Loaded Configuration ================"
	bottomTemplate = "======================================================"
)

func Load(print bool) *Config {
	k := koanf.New(delimeter)

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	if err := loadEnv(k); err != nil {
		log.Printf("error loading environment variable: %s", err)
	}

	config := Config{}
	tag := koanf.UnmarshalConf{Tag: tagName}

	if err := k.UnmarshalWithConf("", &config, tag); err != nil {
		log.Fatalf("error unmarshal config: %v", err)
	}

	if print {
		log.Printf("%s\n%v\n%s", upTemplate, spew.Sdump(config), bottomTemplate)
	}

	return &config
}

func loadEnv(k *koanf.Koanf) error {
	callback := func(source string) string {
		base := strings.ToLower(strings.TrimPrefix(source, envPrefix))
		return strings.ReplaceAll(base, seperator, delimeter)
	}

	if err := k.Load(env.Provider(envPrefix, delimeter, callback), nil); err != nil {
		return fmt.Errorf("error loading environment variable: %s", err)
	}

	return nil
}
