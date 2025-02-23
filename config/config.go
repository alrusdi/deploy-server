package config

type Config struct {
	Shell struct {
		Binary string   `yaml:"binary"`
		Args   []string `yaml:"args"`
	} `yaml:"shell"`
}
