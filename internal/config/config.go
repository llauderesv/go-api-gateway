package config

import (
	"fmt"
	"net/url"
	"strings"
)

type Route struct {
	Path       string `yaml:"path"`
	Target     string `yaml:"target"`
	TargetPath string `yaml:"targetPath"`
}

type Config struct {
	Port   int     `yaml:"port"`
	Routes []Route `yaml:"routes"`
}

func (c Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if len(c.Routes) == 0 {
		return fmt.Errorf("at least one route is required")
	}

	for i, route := range c.Routes {
		if err := validateRoute(route); err != nil {
			return fmt.Errorf("route %d: %w", i, err)
		}
	}

	return nil
}

func validateRoute(route Route) error {
	if strings.TrimSpace(route.Path) == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !strings.HasPrefix(route.Path, "/") {
		return fmt.Errorf("path must start with /")
	}

	if strings.TrimSpace(route.Target) == "" {
		return fmt.Errorf("target cannot be empty")
	}

	target, err := url.Parse(route.Target)
	if err != nil {
		return fmt.Errorf("invalid target: %w", err)
	}

	if target.Scheme != "http" && target.Scheme != "https" {
		return fmt.Errorf("target must use http or https")
	}

	if target.Host == "" {
		return fmt.Errorf("target must contain a host")
	}

	if strings.TrimSpace(route.TargetPath) == "" {
		return fmt.Errorf("targetPath cannot be empty")
	}

	if !strings.HasPrefix(route.TargetPath, "/") {
		return fmt.Errorf("targetPath must start with /")
	}

	return nil
}
