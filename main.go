package main

import (
	"fmt"
	"os"

	"bufio"
	"context"
	"encoding/json"
	"io"

	"github.com/urfave/cli"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/docker/distribution/reference"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "registry,r",
			EnvVar: "PLUGIN_REGISTRY",
		},
		cli.StringSliceFlag{
			Name:   "images,i",
			EnvVar: "PLUGIN_IMAGES",
		},
	}
	app.Action = func(c *cli.Context) error {
		d, err := NewDocker(os.Stderr)
		if err != nil {
			return err
		}
		images := c.StringSlice("images")
		registry := c.String("registry")
		for _, image := range images {
			n, err := reference.ParseNormalizedNamed(image)
			if err != nil {
				return err
			}
			repository := reference.FamiliarName(n)
			tag := getAPITagFromNamedRef(n)
			newImage := fmt.Sprintf("%s/%s:%s", registry, repository, tag)

			if err := d.pull(image); err != nil {
				return err
			}
			if err := d.tag(image, newImage); err != nil {
				return err
			}
			if err := d.push(newImage); err != nil {
				return err
			}
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		os.Exit(1)
	}
}

func getAPITagFromNamedRef(ref reference.Named) string {
	if digested, ok := ref.(reference.Digested); ok {
		return digested.Digest().String()
	}
	ref = reference.TagNameOnly(ref)
	if tagged, ok := ref.(reference.Tagged); ok {
		return tagged.Tag()
	}
	return ""
}

type Docker struct {
	cli    *docker.Client
	stdout io.Writer
}

func NewDocker(stdout io.Writer) (*Docker, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &Docker{cli, stdout}, nil
}

func (d *Docker) pull(image string) error {
	resp, err := d.cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	return d.handleResponse(resp)
}

func (d *Docker) tag(image string, tag string) error {
	return d.cli.ImageTag(context.Background(), image, tag)
}

func (d *Docker) push(image string) error {
	resp, err := d.cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		RegistryAuth: "dummy",
	})
	if err != nil {
		return err
	}
	return d.handleResponse(resp)
}

func (d *Docker) handleResponse(resp io.ReadCloser) error {
	r := bufio.NewReader(resp)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		var v map[string]interface{}
		if err := json.Unmarshal(line, &v); err != nil {
			return fmt.Errorf("failed to unmarshal response. (line=%s, err=%v)", string(line), err)
		}
		status, exist := v["status"]
		if !exist {
			continue
		}
		fmt.Fprintln(d.stdout, status)
	}
}
