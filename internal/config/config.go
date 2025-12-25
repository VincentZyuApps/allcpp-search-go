package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

func Load() *Config {
	cfg := &Config{
		Host:  "0.0.0.0",
		Port:  51225,
		Debug: false,
	}

	// 尝试读取配置文件
	configPath := findConfigFile()
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			log.Printf("⚠️  读取配置文件失败: %v，使用默认配置", err)
		} else {
			log.Printf("✅ 已加载配置文件: %s", configPath)
		}
	} else {
		log.Println("⚠️  未找到 config.yml，使用默认配置")
	}

	// 环境变量可以覆盖配置文件
	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if debug := os.Getenv("DEBUG"); debug == "true" || debug == "1" {
		cfg.Debug = true
	}

	return cfg
}

// findConfigFile 查找配置文件
func findConfigFile() string {
	// 可能的配置文件位置
	possiblePaths := []string{
		"config.yml",
		"config.yaml",
	}

	// 获取可执行文件所在目录
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		possiblePaths = append(possiblePaths,
			filepath.Join(execDir, "config.yml"),
			filepath.Join(execDir, "config.yaml"),
		)
	}

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err == nil && workDir != filepath.Dir(execPath) {
		possiblePaths = append(possiblePaths,
			filepath.Join(workDir, "config.yml"),
			filepath.Join(workDir, "config.yaml"),
		)
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// loadFromFile 从文件加载配置
func loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}
