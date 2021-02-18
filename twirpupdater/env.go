package twirpupdater

import (
	"github.com/thepwagner/action-update/actions/updateaction"
	"github.com/thepwagner/action-update/updater"
)

type Environment struct {
	updateaction.Environment

	// FIXME: remove default default
	ImageName string `env:"INPUT_UPDATER_IMAGE" envDefault:"action-update-twirp-gradle"`
}

func (e *Environment) NewUpdater(root string) updater.Updater {
	return NewUpdater(root, e.ImageName)
}
