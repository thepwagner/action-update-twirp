package twirpupdater

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update-twirp/proto/actionupdate/v1"
	"github.com/thepwagner/action-update/updater"
)

type Updater struct {
	root        string
	updaterName string
	imageName   string

	mu        sync.Mutex
	container *Container
}

func NewUpdater(root, imageName string) *Updater {
	updaterName := fmt.Sprintf("twirp-%s", strings.Split(imageName, ":")[0])
	return &Updater{
		root:        root,
		updaterName: updaterName,
		imageName:   imageName,
	}
}

var _ updater.Updater = (*Updater)(nil)

func (u *Updater) Name() string {
	return u.updaterName
}

func (u *Updater) Dependencies(ctx context.Context) ([]updater.Dependency, error) {
	svc, err := u.updateService(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating update container: %w", err)
	}

	res, err := svc.ListDependencies(ctx, &v1.ListDependenciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("listing dependencies: %w", err)
	}
	deps := res.GetDependencies()

	ret := make([]updater.Dependency, 0, len(deps))
	for _, dep := range deps {
		ret = append(ret, depFromProto(dep))
	}
	return ret, nil
}

func (u *Updater) updateService(ctx context.Context) (v1.UpdateService, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.container == nil {
		logrus.WithField("image", u.imageName).Info("booting updater container")
		ctr, err := newContainer(ctx, u.root, u.imageName)
		if err != nil {
			return nil, err
		}
		u.container = ctr
	}
	return u.container.UpdateService(), nil
}

func (u *Updater) Check(context.Context, updater.Dependency, func(string) bool) (*updater.Update, error) {
	// TODO: implement
	return nil, nil
}

func (u *Updater) ApplyUpdate(context.Context, updater.Update) error {
	// TODO: implement
	return nil
}

func (u *Updater) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.container == nil {
		logrus.Debug("closing updater without container")
		return nil
	}
	return u.container.Close()
}
