package task

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Docker: The Docker type will encapsulate everything we need to run a task in a container.
type Docker struct {
	Client      DockerApi
	Config      Config
	ContainerId string
	ioCopy      func(dst io.Writer, src io.Reader) (written int64, err error)
	stdCopy     func(dstout, dsterr io.Writer, src io.Reader) (written int64, err error)
}

// Config: The Config type represents the configuration for orchestration tasks.
type Config struct {
	AttachStderr  bool
	AttachStdin   bool
	AttachStdout  bool
	Cmd           []string
	Disk          int64
	Env           []string
	Image         string
	Memory        int64
	Name          string
	RestartPolicy string // Accepted values: always, unless-stopped, on-failure and empty.
}

// DockerApi: The DockerApi type represents the interface for the Docker API.
type DockerApi interface {
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, id string, opts types.ContainerStartOptions) error
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
}

// NewDocker: The NewDocker function returns a new Docker.
func NewDocker(client DockerApi) *Docker {
	return &Docker{
		Client: client,
	}
}

// DockerResult: The DockerResult type represents the result of running a task in a container.
type DockerResult struct {
	Action      string
	ContainerId string
	Error       error
	Result      string
}

// Run: The Run function runs a task in a container.
func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	ioReadCloser, err := d.Client.ImagePull(
		ctx,
		d.Config.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	d.ioCopy(os.Stdout, ioReadCloser)
	restartPolicy := container.RestartPolicy{Name: d.Config.RestartPolicy}
	resources := container.Resources{
		Memory: d.Config.Memory,
	}
	config := container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}
	hostConfig := container.HostConfig{
		Resources:       resources,
		RestartPolicy:   restartPolicy,
		PublishAllPorts: true,
	}
	res, err := d.Client.ContainerCreate(
		ctx, &config, &hostConfig, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf(
			"Error creating container using image %s: %v\n",
			d.Config.Image, err,
		)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v\n", res.ID, err)
		return DockerResult{Error: err}
	}
	d.ContainerId = res.ID
	out, err := d.Client.ContainerLogs(
		ctx,
		res.ID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		},
	)
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", res.ID, err)
		return DockerResult{Error: err}
	}
	d.stdCopy(os.Stdout, os.Stderr, out)
	return DockerResult{
		Action:      "start",
		ContainerId: res.ID,
		Result:      "success",
	}
}
