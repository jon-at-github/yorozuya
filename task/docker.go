package task

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
)

// Docker: The Docker type will encapsulate everything we need to run a task in a container.
type Docker struct {
	Clients     *Clients
	Config      Config
	ContainerId string
}

// Clients: The Clients type represents the clients we need to run a task in a container.
type Clients struct {
	DockerApi DockerApi
	IO        IO
}

type IO interface {
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
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
}

// NewDocker: The NewDocker function returns a new Docker.
func NewDocker(clients *Clients) *Docker {
	return &Docker{
		Clients: clients,
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
	ioReadCloser, err := d.Clients.DockerApi.ImagePull(
		context.Background(),
		d.Config.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	d.Clients.IO.Copy(os.Stdout, ioReadCloser)
	return DockerResult{}
}
