package plugin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

// Settings for the Plugin.
type Settings struct {
	DroneToken string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if p.settings.DroneToken == "" {
		return errors.New("mandatory DRONE_TOKEN is not set")
	}
	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	client := newDroneClient("https://"+p.pipeline.System.Host, p.settings.DroneToken)

	// p.pipeline.Build uses a different struct with less information
	currentBuild, err := client.Build(
		p.pipeline.Repo.Owner,
		p.pipeline.Repo.Name,
		p.pipeline.Build.Number,
	)
	if err != nil {
		return err
	}

	builds, err := client.BuildList(
		p.pipeline.Repo.Owner,
		p.pipeline.Repo.Name,
		drone.ListOptions{
			Page: 1,
			Size: 100,
		},
	)
	if err != nil {
		return err
	}

	for _, build := range builds {
		if build.Status == drone.StatusRunning &&
			build.Ref == currentBuild.Ref &&
			build.ID < currentBuild.ID {
			// kill only running builds for the same reference with a lower id
			err := client.BuildCancel(
				p.pipeline.Repo.Owner,
				p.pipeline.Repo.Name,
				int(build.Number),
			)

			if err != nil {
				fmt.Printf("could not cancel build %v because of %v", build.Number, build.Error)
				continue
			}

			fmt.Printf("cancelled build %v", build.ID)
		}
	}

	return err
}

func newDroneClient(droneURL, token string) drone.Client {
	config := new(oauth2.Config)
	auth := config.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)

	auth.CheckRedirect = func(*http.Request, []*http.Request) error {
		return fmt.Errorf("attempting to redirect the requests. Did you configure the correct drone server address?")
	}

	return drone.NewClient(droneURL, auth)
}
