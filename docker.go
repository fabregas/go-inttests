package inttest_utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerContainer struct {
	conf     *container.Config
	hostConf *container.HostConfig
	netConf  *network.NetworkingConfig
	cli      *client.Client
	id       string
}

func NewDockerContainer(image string) (*DockerContainer, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

getimg:
	res, err := cli.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: filters.NewArgs(filters.KeyValuePair{"reference", image}),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		// try pull image
		fmt.Printf("=> pulling image %s\n", image)
		pullout, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}
		defer pullout.Close()
		if err != nil {
			out, _ := ioutil.ReadAll(pullout)
			fmt.Println(string(out))
			return nil, err
		}
		goto getimg

	}
	imageId := res[0].ID

	return &DockerContainer{
		cli:      cli,
		conf:     &container.Config{Image: imageId, Tty: true},
		hostConf: &container.HostConfig{},
		netConf:  &network.NetworkingConfig{},
	}, nil
}

func (c *DockerContainer) Start(name string) error {
	cont, err := c.cli.ContainerCreate(context.Background(), c.conf, c.hostConf, c.netConf, name)
	if err != nil {
		return err
	}

	c.id = cont.ID
	err = c.cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *DockerContainer) SetEnv(env map[string]string) {
	envreq := make([]string, 0, len(env))
	for k, v := range env {
		envreq = append(envreq, fmt.Sprintf("%s=%s", k, v))
	}

	c.conf.Env = envreq
}

func (c *DockerContainer) SetPortsBinding(portsBinds map[string]string) {
	c.hostConf.PortBindings = make(nat.PortMap)
	for cport, v := range portsBinds {
		vp := strings.Split(v, ":")
		var host, port string
		if len(vp) == 1 {
			port = vp[0]
		} else {
			host = vp[0]
			port = vp[1]
		}
		c.hostConf.PortBindings[nat.Port(fmt.Sprintf("%s/tcp", cport))] = []nat.PortBinding{nat.PortBinding{host, port}}
	}
}

func (c *DockerContainer) IP() (string, error) {
	ret, err := c.cli.ContainerInspect(context.Background(), c.id)
	if err != nil {
		return "", err
	}
	return ret.NetworkSettings.IPAddress, nil
}

func (c *DockerContainer) Stop() error {
	err := c.cli.ContainerKill(context.Background(), c.id, "SIGINT")
	if err != nil {
		return err
	}

	return c.cli.ContainerRemove(
		context.Background(),
		c.id,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		},
	)
}

func (c *DockerContainer) Logs() (string, error) {
	reader, err := c.cli.ContainerLogs(
		context.Background(),
		c.id,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		},
	)
	if err != nil {
		return "", err
	}

	out, err := ioutil.ReadAll(reader)
	reader.Close()
	return string(out), err
}
