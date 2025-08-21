package models

import "time"

// Config represents the configuration for the application
type Config struct {
	Behavior Behavior `mapstructure:"behavior" validate:"required"`

	Activities []Activity `mapstructure:"activities" validate:"required"`

	Debug   bool   `mapstructure:"debug"`
	LogFile string `mapstructure:"logfile"`
}

type Behavior struct {
	// IdleTimeout is the duration of user inactivity after which the activities should be triggered.
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`

	// ResumeAfterInactivity indicates whether to resume activities after a period of inactivity.
	ResumeAfterInactivity bool `mapstructure:"resume_after_inactivity"`

	// PauseWhenUserIsActive indicates whether to pause activities when the user is active.
	PauseWhenUserIsActive bool `mapstructure:"pause_when_user_is_active"`
}

type Schedule struct {
	Enabled bool           `mapstructure:"enabled"`
	From    string         `mapstructure:"from"`
	To      string         `mapstructure:"to"`
	Days    []time.Weekday `mapstructure:"days"`
}

type Activity struct {
	Kind     Kind          `mapstructure:"kind" validate:"required"`
	Enabled  bool          `mapstructure:"enabled"`
	Schedule Schedule      `mapstructure:"schedule"`
	Interval time.Duration `mapstructure:"interval" validate:"required,gte=5000000000"`
}
