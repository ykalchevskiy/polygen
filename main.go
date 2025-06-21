package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configPath := flag.String("config", ".polygen.json", "Path to the configuration file")
	flag.Parse()

	if err := run(*configPath); err != nil {
		log.Fatalf("Failed to generate: %v", err)
	}
}

func run(configPath string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file '%s': %v", configPath, err)
	}

	var config FileConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("parsing config file '%s': %v", configPath, err)
	}

	configDir := filepath.Dir(configPath)

	for _, typeConfig := range config.Types {
		cfg := convertFileConfigToConfig(&typeConfig, &config)

		outputPath := getOutputPath(&typeConfig, configDir)

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("creating output directory '%s' for type '%s': %v", outputPath, typeConfig.Type, err)
		}

		code, err := generate(cfg)
		if err != nil {
			return fmt.Errorf("generating code for type '%s': %v", typeConfig.Type, err)
		}

		if err := os.WriteFile(outputPath, code, 0644); err != nil {
			return fmt.Errorf("writing generated code for type '%s': %v", typeConfig.Type, err)
		}
	}

	return nil
}
