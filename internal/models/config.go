package models

// Config represents the configuration for the application
type Config struct {
	Interval int    `yaml:"interval"`
	Distance int    `yaml:"distance"`
	Debug    bool   `yaml:"debug"`
	LogFile  string `yaml:"logfile"`
}
