package task

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type DockerFixture struct {
	*gunit.Fixture

	docker         *Docker
	fakeDockerApi  *FakeDockerApi
	fakeIO         *FakeIO
	FakeReadCloser *FakeReadCloser
	fakeReader     *FakeReader
	FakeWriter     *FakeWriter
}

type FakeDockerApi struct {
	err     error
	options *types.ImagePullOptions
	output  io.ReadCloser
	refStr  string
}

func (f *FakeDockerApi) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	f.refStr = refStr
	f.options = &options
	return f.output, f.err
}

type FakeIO struct {
	dst io.Writer
	src io.Reader
}

type FakeWriter struct{}

func (f *FakeWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type FakeReader struct{}

func (f *FakeReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f *FakeIO) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	f.dst = dst
	f.src = src
	return 0, nil
}

func (f *DockerFixture) Setup() {
	f.fakeDockerApi = new(FakeDockerApi)
	f.FakeReadCloser = new(FakeReadCloser)
	f.fakeDockerApi.output = f.FakeReadCloser
	f.fakeReader = new(FakeReader)
	f.FakeWriter = new(FakeWriter)
	f.fakeIO = &FakeIO{
		dst: f.FakeWriter,
		src: f.fakeReader,
	}
	f.docker = NewDocker(
		&Clients{
			DockerApi: f.fakeDockerApi,
			IO:        f.fakeIO,
		},
	)
}

type FakeReadCloser struct{}

func (f *FakeReadCloser) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (f *FakeReadCloser) Close() error {
	return nil
}

func (f *DockerFixture) TestRunPullsImageAndWritesToStdout() {
	f.docker.Config.Image = "test-image"
	result := f.docker.Run()
	f.So(result.Error, should.BeNil)
	f.AssertEqual("test-image", f.fakeDockerApi.refStr)
	f.AssertEqual(f.fakeDockerApi.output, f.fakeIO.src)
	f.So(f.fakeIO.src, should.Resemble, f.fakeDockerApi.output)
}

func TestDockerFixture(t *testing.T) {
	gunit.Run(new(DockerFixture), t)
}
