package cfg

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type BuildConfig struct {
	ConfigVersion string         `yaml:"version"`
	Setup         []BuildStep    `yaml:"setup"`
	Steps         []BuildStep    `yaml:"steps"`
	PostBuild     []BuildStep    `yaml:"post_build"`
	Parallel      []ParallelStep `yaml:"parallel"`
	Ignore        []string       `yaml:"ignore"`
}

type BuildStep struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	If      string `yaml:"if,omitempty"`
}

type ParallelStep struct {
	Name     string   `yaml:"name"`
	Commands []string `yaml:"commands"`
}

func LoadConfig(path string) (*BuildConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config BuildConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *BuildConfig) ShouldIgnore(path string) bool {
	base := filepath.Base(path)

	if strings.HasPrefix(base, ".") {
		return true
	}

	for _, pattern := range c.Ignore {
		matched, _ := filepath.Match(pattern, base)
		if matched || strings.HasPrefix(path, pattern) {
			return true
		}
	}
	return false
}
