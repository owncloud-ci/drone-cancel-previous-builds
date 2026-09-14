package main

import (
	"github.com/owncloud-ci/drone-cancel-previous-builds/plugin"
	"github.com/urfave/cli/v3"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "drone-token",
			Value:       "",
			EnvVars:     []string{"DRONE_TOKEN", "PLUGIN_DRONE_TOKEN"},
			Usage:       "Drone token",
			Destination: &settings.DroneToken,
		},
	}
}
