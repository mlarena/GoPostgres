package config

import (
	"fmt"
	"os"
	"time"
	
	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type MonitoringConfig struct {
	LongQueryThreshold string `yaml:"long_query_threshold"`
	CheckInterval      string `yaml:"check_interval"`
	MaxDatabases       int    `yaml:"max_databases"`
}

type LoggingConfig struct {
	Level     string `yaml:"level"`
	Path      string `yaml:"path"`
	MaxSize   int    `yaml:"max_size"`    // в мегабайтах
	MaxBackups int   `yaml:"max_backups"` // количество файлов
	MaxAge    int    `yaml:"max_age"`     // в днях
}

type Config struct {
	Database    map[string]DatabaseConfig `yaml:"database"`
	Monitoring  MonitoringConfig         `yaml:"monitoring"`
	Logging     LoggingConfig            `yaml:"logging"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}
	
	return cfg, nil
}

func (d DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

func (m MonitoringConfig) GetCheckInterval() (time.Duration, error) {
	return time.ParseDuration(m.CheckInterval)
}

func (m MonitoringConfig) GetLongQueryThreshold() (time.Duration, error) {
	return time.ParseDuration(m.LongQueryThreshold)
}