package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Config represents the configuration for the application
type Config struct {
	Behavior Behavior `mapstructure:"behavior" validate:"required"`

	Activities []Activity `mapstructure:"activities" validate:"required,dive"`

	Debug   bool   `mapstructure:"debug"`
	LogFile string `mapstructure:"logfile"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Behavior, validation.Required),
		validation.Field(&c.Activities),
		validation.Field(&c.Debug),
		validation.Field(&c.LogFile),
	)
}

// Behavior defines the behavior settings for the application.
type Behavior struct {
	// StartActivitiesOnStartup indicates whether the application should start activities on startup.
	StartActivitiesOnStartup bool `mapstructure:"start_on_boot"`

	// IdleTimeout is the duration of user inactivity after which the activities should be triggered.
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`

	// ResumeAfterInactivity indicates whether to resume activities after a period of inactivity.
	ResumeAfterInactivity bool `mapstructure:"resume_after_inactivity"`

	// PauseWhenUserIsActive indicates whether to pause activities when the user is active.
	PauseWhenUserIsActive bool `mapstructure:"pause_when_user_is_active"`
}

// Schedule defines the schedule settings for the application.
type Schedule struct {
	Enabled *bool          `mapstructure:"enabled"`
	From    string         `mapstructure:"from"`
	To      string         `mapstructure:"to"`
	Days    []time.Weekday `mapstructure:"days"`
}

// Activity defines the activity settings for the application.
type Activity struct {
	Kind     Kind          `mapstructure:"kind"`
	Pattern  Pattern       `mapstructure:"pattern"`
	Enabled  *bool         `mapstructure:"enabled"`
	Schedule Schedule      `mapstructure:"schedule"`
	Interval time.Duration `mapstructure:"interval"`
}

func (a Activity) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Kind, validation.Required, validation.In(KindMouse, KindKeyboard)),
		validation.Field(&a.Pattern, validation.Required.When(a.Kind == KindMouse), validation.In(PatternSquare, PatternTriangle, PatternUpAndDown, PatternLeftAndRight)),
		validation.Field(&a.Enabled, validation.NotNil),
		validation.Field(&a.Interval, validation.Required, validation.Min(5*time.Second)),
	)
}
