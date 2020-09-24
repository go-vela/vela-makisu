// Copyright (c) 2020 Target Brands, Inr. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/uber/makisu/lib/registry"
	"github.com/urfave/cli/v2"
)

const (
	// dockerHubConf represents the config that provides authentication to docker.index.io.
	dockerHubConf = `{
		"index.docker.io": {
      ".*": {
         "security": {
            "tls": {
               "client": {
                  "disabled": false
               }
            },
            "basic": {
               "username": "",
               "password": ""
            }
         }
      }
   }		
 }`

	// mirrorConf represents the config that provides authentication to a registry mirror.
	mirrorConf = `{
		"%s": {
      ".*": {
         "security": {
            "tls": {
               "client": {
                  "disabled": false
               }
            },
            "basic": {
               "username": "",
               "password": ""
            }
         }
      }
   }		
 }`

	// registryConf represents the config that provides authentication to a registry.
	registryConf = `{
		"%s": {
      ".*": {
         "security": {
            "tls": {
               "client": {
                  "disabled": false
               }
            },
            "basic": {
               "username": "%s",
               "password": "%s"
            }
         }
      }
   }		
 }`
)

// Registry represents the input parameters for the plugin.
type Registry struct {
	// full url to a Docker Registry mirror
	Mirror string
	// full url to Docker Registry
	Name string
	// password for communication with the Docker Registry
	Password string
	// user name for communication with the Docker Registry
	Username string
}

var (
	// appFs represents a instance of the filesystem.
	appFS = afero.NewOsFs()

	// configFlags represents for config settings on the cli.
	// nolint
	configFlags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "REGISTRY_NAME"},
			FilePath: string("/vela/parameters/makisu/registry/name,/vela/secrets/docker/registry/name"),
			Name:     "registry.name",
			Usage:    "Docker registry address to communicate with",
			Value:    "index.docker.io",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_MIRROR", "REGISTRY_MIRROR"},
			FilePath: string("/vela/parameters/makisu/registry/name,/vela/secrets/docker/registry/name"),
			Name:     "registry.mirror",
			Usage:    "Docker registry mirror address to communicate with",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "REGISTRY_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: string("/vela/parameters/makisu/registry/password,/vela/secrets/makisu/registry/password,/vela/secrets/makisu/password"),
			Name:     "registry.password",
			Usage:    "password for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "REGISTRY_USERNAME", "DOCKER_USERNAME"},
			FilePath: string("/vela/parameters/makisu/registry/username,/vela/secrets/makisu/registry/username,/vela/secrets/makisu/username"),
			Name:     "registry.username",
			Usage:    "user name for communication with the registry",
		},
	}

	// configPath represents the location of the Docker config file for setting registries.
	configPath = "/makisu/registry/config.json"
)

// Write creates a Docker config.json file for building and publishing the image.
func (r *Registry) Write() error {
	logrus.Trace("creating registry configuration file")

	// allocate a config registry map
	config := make(registry.Map)

	// add the anonymous docker hub config to map
	err := json.Unmarshal([]byte(dockerHubConf), &config)
	if err != nil {
		return err
	}

	// when a mirror is provided add it to the config
	if len(r.Mirror) != 0 {
		// create output string for config.json file
		mirror := fmt.Sprintf(
			mirrorConf,
			r.Mirror,
		)

		// add the user config to the registry map
		err = json.Unmarshal([]byte(mirror), &config)
		if err != nil {
			return err
		}
	}

	// create output string for config.json file
	registry := fmt.Sprintf(
		registryConf,
		r.Name,
		r.Username,
		r.Password,
	)

	// add the user config to the registry map
	err = json.Unmarshal([]byte(registry), &config)
	if err != nil {
		return err
	}

	registryConf, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	return a.WriteFile(configPath, registryConf, 0644)
}

// Validate verifies the Config is properly configured.
func (r *Registry) Validate() error {
	logrus.Trace("validating config plugin configuration")

	// verify password are provided
	if len(r.Password) == 0 {
		return fmt.Errorf("no config password provided")
	}

	// verify url is provided
	if len(r.Name) == 0 {
		return fmt.Errorf("no config address provided")
	}

	// verify username is provided
	if len(r.Username) == 0 {
		return fmt.Errorf("no config username provided")
	}

	return nil
}
