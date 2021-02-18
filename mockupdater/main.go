package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update-twirp/mockupdater/updater"
	v1 "github.com/thepwagner/action-update-twirp/proto/actionupdate/v1"
)

func main() {
	svc := updater.NewUpdateService()
	handler := v1.NewUpdateServiceServer(svc)

	port := int64(9600)
	if portEnv, ok := os.LookupEnv("ACTION_UPDATE_TWIRP_PORT"); ok {
		if parsedPort, err := strconv.ParseInt(portEnv, 10, 64); err == nil {
			port = parsedPort
		}
	}

	mux := http.NewServeMux()
	mux.Handle(handler.PathPrefix(), handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatalf("http error")
	}
}
