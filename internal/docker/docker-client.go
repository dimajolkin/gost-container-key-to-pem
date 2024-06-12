package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
)

type Status string

const (
	StatusCreate Status = "create"
	StatusFailed Status = "failed"
	StatusOk     Status = "ok"
)

type Client struct {
	Status  Status
	cli     *client.Client
	context context.Context
}

func NewDockerClient() Client {
	return Client{
		Status:  StatusCreate,
		context: context.Background(),
	}
}

func (obj *Client) Connect() error {
	var err error
	obj.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		obj.Status = StatusFailed
		return err
	}
	defer obj.cli.Close()
	obj.Status = StatusOk

	return nil
}

func (obj Client) PullImage(imageName string) (*io.ReadCloser, error) {
	if obj.Status != StatusOk {
		return nil, errors.New(fmt.Sprintf("failed to pull image. Status: %s", obj.Status))
	}

	reader, err := obj.cli.ImagePull(obj.context, imageName, image.PullOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return &reader, nil
}

func (obj Client) Run(imageName string, name string, cmd []string) (string, error) {
	if obj.Status != StatusOk {
		return "", errors.New(fmt.Sprintf("failed to run image. Status: %s", obj.Status))
	}

	err := obj.removeContainer(name)
	if err != nil {
	}
	resp, err := obj.cli.ContainerCreate(obj.context, &container.Config{
		Image: imageName,
		Cmd:   cmd,
	}, nil, nil, nil, name)
	if err != nil {
		return "", err
	}

	if err := obj.cli.ContainerStart(obj.context, name, container.StartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (obj Client) Wait(name string) (*io.ReadCloser, error) {
	statusCh, errCh := obj.cli.ContainerWait(obj.context, name, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := obj.cli.ContainerLogs(obj.context, name, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (obj Client) removeContainer(name string) error {
	err := obj.cli.ContainerRemove(obj.context, name, container.RemoveOptions{Force: true})
	if err != nil {
		//if err.Error() == "Error response from daemon: No such container: "+name {
		//	return nil
		//}
		//
		//if err.Error() == "Error response from daemon: No such container: "+name {
		//	return nil
		//}

		return err
	}

	statusCh, errCh := obj.cli.ContainerWait(obj.context, name, container.WaitConditionRemoved)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

func (obj Client) Copy(name string, distPath string, content io.Reader) error {
	return obj.cli.CopyToContainer(obj.context, name, distPath, content, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                true,
	})
}

func (obj Client) Exec(name string, cmd []string) error {
	response, err := obj.cli.ContainerExecCreate(obj.context, name, types.ExecConfig{Cmd: cmd, Tty: true, AttachStdin: true, AttachStderr: true})
	if err != nil {
		return err
	}

	return obj.cli.ContainerExecStart(obj.context, response.ID, types.ExecStartCheck{})
}

func (obj Client) Test() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, "sokolko/export-cryptopro-cert", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "sokolko/export-cryptopro-cert",

		Cmd: []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
