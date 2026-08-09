package config

type Route struct {
	Path       string `yaml:"path"`
	Target     string `yaml:"target"`
	TargetPath string `yaml:"targetPath"`
}

type Config struct {
	Port   int     `yaml:"port"`
	Routes []Route `yaml:"routes"`
}
