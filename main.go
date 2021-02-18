package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update-twirp/twirpupdater"
	"github.com/thepwagner/action-update/actions/updateaction"
)

func main() {
	var cfg twirpupdater.Environment
	handlers := updateaction.NewHandlers(&cfg)
	ctx := context.Background()
	if err := handlers.ParseAndHandle(ctx, &cfg); err != nil {
		logrus.WithError(err).Fatal("failed")
	}
}
