package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

type BuildConfig struct {
	Inherit       string            `yaml:"inherit"`
	ConfigVersion string            `yaml:"version"`
	Setup         []BuildStep       `yaml:"setup,omitempty"`
	Steps         []BuildStep       `yaml:"steps,omitempty"`
	PostBuild     []BuildStep       `yaml:"post_build,omitempty"`
	Parallel      []ParallelStep    `yaml:"parallel,omitempty"`
	Ignore        []string          `yaml:"ignore,omitempty"`
	Variables     map[string]string `yaml:"vars,omitempty"`
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

	if config.Inherit != "" {
		basePath := config.Inherit
		if !filepath.IsAbs(basePath) {
			basePath = filepath.Join(filepath.Dir(path), basePath)
		}

		baseConfig, err := LoadConfig(basePath)
		if err != nil {
			return nil, fmt.Errorf("❌ Failed to load inherited config: %w", err)
		}

		mergeConfigs(baseConfig, &config)
	}

	// Load .gitignore from CWD
	gitignorePath := filepath.Join(filepath.Dir(path), ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		gitignoreLines, err := os.ReadFile(gitignorePath)
		if err != nil {
			fmt.Println("⚠️ Failed to read .gitignore: ", err)
		} else {
			lines := strings.Split(string(gitignoreLines), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				config.Ignore = append(config.Ignore, line)
			}
		}
	}

	return &config, nil
}

func (c *BuildConfig) ShouldIgnore(path string) bool {
	normalizedPath := filepath.ToSlash(path)

	for _, pattern := range c.Ignore {
		normalizedPattern := filepath.ToSlash(pattern)

		if matched, _ := doublestar.Match(normalizedPattern, normalizedPath); matched {
			return true
		}

		if strings.HasSuffix(normalizedPattern, "/") {
			if matched, _ := doublestar.Match(normalizedPattern+"**", normalizedPath); matched {
				return true
			}
		}
	}
	return false
}

func mergeConfigs(base, override *BuildConfig) {
	if base.Variables != nil {
		if override.Variables == nil {
			override.Variables = make(map[string]string)
		}
		for k, v := range base.Variables {
			if _, exists := override.Variables[k]; !exists {
				override.Variables[k] = v
			}
		}
	}

	override.Ignore = append(base.Ignore, override.Ignore...)
	override.Setup = append(base.Setup, override.Setup...)
	override.Steps = append(base.Steps, override.Steps...)
	override.PostBuild = append(base.PostBuild, override.PostBuild...)
	override.Parallel = append(base.Parallel, override.Parallel...)
}
