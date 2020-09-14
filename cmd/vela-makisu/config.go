// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/base64"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	credentials = `%s:%s`

	registryFile = `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`
)

// Config holds input parameters for the plugin.
type Config struct {
	// enables building images without publishing to the registry
	DryRun bool
	// password for communication with the Docker Registry
	Password string
	// config path the docker json file exists for authentication
	Path string
	// full url to Docker Registry
	URL string
	// user name for communication with the Docker Registry
	Username string
}

var (

	// configFlags represents for config settings on the cli.
	configFlags = []cli.Flag{
		// nolint
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DRY_RUN", "REGISTRY_DRY_RUN"},
			FilePath: string("/vela/parameters/makisu/registry/path,/vela/secrets/makisu/registry/path,/vela/secrets/makisu/path"),
			Name:     "config.dry-run",
			Usage:    "enables building images without publishing to the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "REGISTRY_NAME"},
			FilePath: string("/vela/parameters/makisu/registry/name,/vela/secrets/docker/registry/name"),
			Name:     "config.registry",
			Usage:    "Docker registry name to communicate with",
			Value:    "index.docker.io",
		},
		// nolint
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "REGISTRY_USERNAME", "DOCKER_USERNAME"},
			FilePath: string("/vela/parameters/makisu/registry/username,/vela/secrets/makisu/registry/username,/vela/secrets/makisu/username"),
			Name:     "config.username",
			Usage:    "user name for communication with the registry",
		},
		// nolint
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "REGISTRY_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: string("/vela/parameters/makisu/registry/password,/vela/secrets/makisu/registry/password,/vela/secrets/makisu/password"),
			Name:     "config.password",
			Usage:    "password for communication with the registry",
		},
	}
)

// Create creates a Docker config.json file for usage with makisu config flag.
func (c *Config) Create() (string, error) {
	logrus.Trace("creating registry configuration file")

	// check if name, username and password are provided
	if len(c.URL) == 0 || len(c.Username) == 0 || len(c.Password) == 0 {
		return "", fmt.Errorf("missing one of the required URL, username, or password parameters")
	}

	// create basic authentication string for config.json file
	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(credentials, c.Username, c.Password)),
	)

	// create output string for config.json file
	return fmt.Sprintf(
		registryFile,
		c.URL,
		basicAuth,
	), nil
}

// Validate verifies the Config is properly configured.
func (c *Config) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify password are provided
	if len(c.Password) == 0 {
		return fmt.Errorf("no config password provided")
	}

	// verify url is provided
	if len(c.URL) == 0 {
		return fmt.Errorf("no config url provided")
	}

	// verify username is provided
	if len(c.Username) == 0 {
		return fmt.Errorf("no config username provided")
	}

	return nil
}
