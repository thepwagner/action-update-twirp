package twirpupdater

import (
	"github.com/thepwagner/action-update-twirp/proto/actionupdate/v1"
	"github.com/thepwagner/action-update/updater"
)

func depFromProto(dep *v1.Dependency) updater.Dependency {
	return updater.Dependency{
		Path:     dep.Path,
		Version:  dep.Version,
		Indirect: dep.Indirect,
	}
}
