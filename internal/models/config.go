package models

import "time"

// Config represents the configuration for the application
type Config struct {
	Interval int `yaml:"interval"`
	Distance int `yaml:"distance"`

	ResumeAfterInactivity      bool `yaml:"resume_after_inactivity"`
	ResumeAfterInactivityValue int  `yaml:"resume_after_inactivity_value"`

	PauseWhenUserIsActive bool `yaml:"pause_when_user_is_active"`

	Schedule Schedule `yaml:"schedule"`

	Debug   bool   `yaml:"debug"`
	LogFile string `yaml:"logfile"`
}

type Schedule struct {
	Enabled bool           `yaml:"enabled"`
	From    string         `yaml:"from"`
	To      string         `yaml:"to"`
	Days    []time.Weekday `yaml:"days"`
}
