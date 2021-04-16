module github.com/thepwagner/action-update-twirp

go 1.15

require (
	github.com/containerd/containerd v1.5.0-rc.1 // indirect
	github.com/docker/docker v20.10.3+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/thepwagner/action-update v0.0.40
	github.com/twitchtv/twirp v7.2.0+incompatible
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/protobuf v1.26.0
)

replace github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
