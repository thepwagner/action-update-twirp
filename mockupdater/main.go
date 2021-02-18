package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update-twirp/mockupdater/updater"
	v1 "github.com/thepwagner/action-update-twirp/proto/actionupdate/v1"
)

func main() {
	svc := updater.NewUpdateService()
	handler := v1.NewUpdateServiceServer(svc)

	mux := http.NewServeMux()
	mux.Handle(handler.PathPrefix(), handler)
	if err := http.ListenAndServe(":9600", handler); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatalf("http error")
	}
}
