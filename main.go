package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configPath := flag.String("config", ".polygen.json", "Path to the configuration file")
	flag.Parse()

	configData, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file '%s': %v", *configPath, err)
	}

	var config FileConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config file '%s': %v", *configPath, err)
	}

	configDir := filepath.Dir(*configPath)

	for _, typeConfig := range config.Types {
		cfg := convertFileConfigToConfig(&typeConfig, &config)

		outputPath := getOutputPath(&typeConfig, configDir)

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			log.Fatalf("Failed to create output directory '%s' for type '%s': %v", outputPath, typeConfig.Type, err)
		}

		code, err := generate(cfg)
		if err != nil {
			log.Fatalf("Failed to generate code for type '%s': %v", typeConfig.Type, err)
		}

		if err := os.WriteFile(outputPath, code, 0644); err != nil {
			log.Fatalf("Failed to write generated code for type '%s': %v", typeConfig.Type, err)
		}
	}
}
