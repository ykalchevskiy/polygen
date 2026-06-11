package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

		switch cfg.JSONVersion {
		case JSONVersionBoth:
			if err := generateAndWrite(cfg, generate, outputPath); err != nil {
				return fmt.Errorf("v1: %v", err)
			}
			outputPathV2 := strings.TrimSuffix(outputPath, ".go") + "_jsonv2.go"
			if err := generateAndWrite(cfg, generateJSONV2, outputPathV2); err != nil {
				return fmt.Errorf("v2: %v", err)
			}
		case JSONVersionV2:
			if err := generateAndWrite(cfg, generateJSONV2, outputPath); err != nil {
				return fmt.Errorf("v2: %v", err)
			}
		default: // JSONVersionV1 or fallback
			if err := generateAndWrite(cfg, generate, outputPath); err != nil {
				return fmt.Errorf("v1: %v", err)
			}
		}
	}

	return nil
}

func generateAndWrite(cfg *Config, gen func(*Config) ([]byte, error), outputPath string) error {
	code, err := gen(cfg)
	if err != nil {
		return fmt.Errorf("generating code for type '%s': %v", cfg.Type, err)
	}
	if err := os.WriteFile(outputPath, code, 0644); err != nil {
		return fmt.Errorf("writing generated code to '%s': %v", outputPath, err)
	}
	return nil
}
