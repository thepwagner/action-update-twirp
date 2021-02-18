package twirpupdater

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	v1 "github.com/thepwagner/action-update-twirp/proto/actionupdate/v1"
)

const (
	apiPort = 9999
	srcPath = "/src"
)

var natPort = nat.Port(fmt.Sprintf("%d/tcp", apiPort))

type Container struct {
	docker  *client.Client
	id      string
	logs    *bytes.Buffer
	apiAddr string
}

func newContainer(ctx context.Context, sourceDir, image string) (*Container, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("connecting to docker: %w", err)
	}

	containerID, err := createContainer(ctx, docker, sourceDir, image)
	if err != nil {
		_ = docker.Close()
		return nil, err
	}

	apiAddr, err := inspectBoundPort(ctx, docker, containerID)
	if err != nil {
		_ = killContainer(docker, containerID)
		return nil, err
	}
	buf, err := setupLogBuffer(ctx, docker, containerID)
	if err != nil {
		_ = killContainer(docker, containerID)
		return nil, err
	}

	if err := waitForAPI(ctx, apiAddr); err != nil {
		_ = killContainer(docker, containerID)
		return nil, err
	}

	return &Container{
		docker:  docker,
		id:      containerID,
		logs:    buf,
		apiAddr: apiAddr,
	}, nil
}

func createContainer(ctx context.Context, docker *client.Client, sourceDir string, image string) (string, error) {
	c := container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			natPort: {},
		},
		Env: []string{
			fmt.Sprintf("ACTION_UPDATE_TWIRP_PORT=%d", apiPort),
			fmt.Sprintf("ACTION_UPDATE_TWIRP_PATH=%s", srcPath),
		},
	}
	h := container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{Type: mount.TypeBind, Source: sourceDir, Target: srcPath},
		},
		PortBindings: nat.PortMap{
			natPort: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "0"}},
		},
	}
	n := network.NetworkingConfig{}
	ctr, err := docker.ContainerCreate(ctx, &c, &h, &n, "")
	if err != nil {
		return "", fmt.Errorf("creating container: %w", err)
	}
	containerID := ctr.ID
	logrus.WithField("container_id", containerID).Debug("created container")

	if err := docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("starting container: %w", err)
	}

	return containerID, nil
}

func inspectBoundPort(ctx context.Context, docker *client.Client, containerID string) (string, error) {
	inspect, err := docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("inspecting container: %w", err)
	}
	apiPortBindings, ok := inspect.NetworkSettings.Ports[natPort]
	if !ok || len(apiPortBindings) == 0 {
		return "", errors.New("api port not bound")
	}
	apiPortBinding := apiPortBindings[0]
	apiAddr := fmt.Sprintf("http://%s:%s", apiPortBinding.HostIP, apiPortBinding.HostPort)
	logrus.WithFields(logrus.Fields{
		"host": apiPortBinding.HostIP,
		"port": apiPortBinding.HostPort,
	}).Debug("bound container to port")
	return apiAddr, nil
}

func setupLogBuffer(ctx context.Context, docker *client.Client, containerID string) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	logs, err := docker.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		Follow:     true,
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		return nil, fmt.Errorf("tailing container logs: %w", err)
	}
	go func() {
		defer logs.Close()
		_, _ = stdcopy.StdCopy(buf, buf, logs)
	}()
	return buf, nil
}

func waitForAPI(ctx context.Context, addr string) error {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	start := time.Now()
	for {
		if time.Since(start) > 30*time.Second {
			return fmt.Errorf("timed out waiting for API")
		}
		if _, err := http.DefaultClient.Do(req); err == nil {
			logrus.WithField("dur", time.Since(start).Milliseconds()).Debug("server ready")
			// any response is fine
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Container) UpdateService() v1.UpdateService {
	return v1.NewUpdateServiceJSONClient(c.apiAddr, http.DefaultClient)
}

func (c *Container) Close() error {
	logrus.WithField("container_id", c.id).Debug("cleaning up container")
	if err := killContainer(c.docker, c.id); err != nil {
		return err
	}
	return nil
}

func killContainer(docker *client.Client, containerID string) error {
	defer docker.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := docker.ContainerKill(ctx, containerID, "SIGKILL"); err != nil {
		return err
	}
	return nil
}
