package start

import (
	"gomove/internal/models"
	"gomove/pkg/activity"
	"gomove/pkg/log"
	"gomove/pkg/watcher"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run(cmd *cobra.Command, args []string) {
	log.Info("gomove is starting", zap.Bool("debug", viper.GetBool("debug")), zap.String("config", viper.ConfigFileUsed()))
	log.Info("press ctrl+c to stop...")

	var config models.Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("unable to decode into struct", zap.Error(err))
	}

	// Validate the configuration
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		log.Fatal("config is not valid", zap.Error(err))
	}

	// Create a new watcher
	watcher := watcher.NewWatcher()

	// Create the activity manager
	activityManager := activity.NewActivityManager(config.Behavior, config.Activities, watcher)

	// Start
	activityManager.Start()
}
