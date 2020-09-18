// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
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

type (
	// Config represents the input parameters for the plugin.
	Config struct {
		// enables building images without publishing to the registry
		DryRun bool
		// Used for translating the raw docker configuration
		Global *Global
		// enables setting configuration for the global flags
		GlobalRaw string
		// password for communication with the Docker Registry
		Password string
		// full url to Docker Registry
		URL string
		// user name for communication with the Docker Registry
		Username string
	}

	// Global represents the global flags that can be set on makisu commands.
	Global struct {
		CPU *CPU
		Log *Log
	}

	// CPU represents the cpu specific global flags available.
	CPU struct {
		// enables viewing the cpu profile of the application
		Profile bool
	}

	// Log represents the log specific global flags available.
	Log struct {
		// enables setting the format output for the logs - options: (json|console)
		Fmt string
		// enables setting the log level for the logs - options: (debug|info|warn|warn)
		Level string
		// enables setting the output path for the logs
		Output string
	}
)

var (
	// appFs represents a instance of the filesystem.
	appFS = afero.NewOsFs()

	// configFlags represents for config settings on the cli.
	// nolint
	configFlags = []cli.Flag{
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
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "REGISTRY_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: string("/vela/parameters/makisu/registry/password,/vela/secrets/makisu/registry/password,/vela/secrets/makisu/password"),
			Name:     "config.password",
			Usage:    "password for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "REGISTRY_USERNAME", "DOCKER_USERNAME"},
			FilePath: string("/vela/parameters/makisu/registry/username,/vela/secrets/makisu/registry/username,/vela/secrets/makisu/username"),
			Name:     "config.username",
			Usage:    "user name for communication with the registry",
		},
	}

	// configPath represents the location of the Docker config file for setting registries.
	configPath = "/makisu/.docker/config.json"

	// configFlags represents the global flags that can be set on the makisu commands.
	// nolint
	globalFlags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_GLOBAL_FLAGS"},
			FilePath: string("/vela/parameters/makisu/config/cpu_profile,/vela/secrets/makisu/config/cpu_profile"),
			Name:     "config.global-flags",
			Usage:    "enables setting the global flags on the CLI",
			Value:    globalDefaultValue,
		},
	}

	// globalDefaultValue represents the default setting for the global config.
	globalDefaultValue = `{
		"cpu":{},		
		"log":{
			 "fmt":"console"
		}
 }`
)

// Write creates a Docker config.json file for building and publishing the image.
func (c *Config) Write() error {
	logrus.Trace("creating registry configuration file")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// create basic authentication string for config.json file
	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(credentials, c.Username, c.Password)),
	)

	// create output string for config.json file
	out := fmt.Sprintf(
		registryFile,
		c.URL,
		basicAuth,
	)

	return a.WriteFile(configPath, []byte(out), 0644)
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (c *Config) Unmarshal() error {
	logrus.Trace("unmarshaling config global flags")

	// allocate configuration to structs
	c.Global = &Global{}

	// check if any global flags were passed
	if len(c.GlobalRaw) > 0 {
		// cast raw global flags into bytes
		globalFlags := []byte(c.GlobalRaw)

		// serialize raw global flags into expected Global type
		err := json.Unmarshal(globalFlags, &c.Global)
		if err != nil {
			return err
		}
	}

	return nil
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

// Flags formats and outputs the flags for
// configuring a Docker.
func (g *Global) Flags() []string {
	logrus.Trace("creating global flags for command")

	// variable to store flags for command
	var flags []string

	// check if cpu profile is provided
	if g.CPU.Profile {
		// add flag for cpu profile from provided build command
		flags = append(flags, "--cpu-profile ")
	}

	// check if log fmt is provided
	if len(g.Log.Fmt) > 0 {
		// add flag for log fmt from provided build command
		flags = append(flags, fmt.Sprintf("--log-fmt=%s", g.Log.Fmt))
	}

	// check if log level is provided
	if len(g.Log.Level) > 0 {
		// add flag for log level from provided build command
		flags = append(flags, fmt.Sprintf("--log-level=%s", g.Log.Level))
	}

	// check if log output is provided
	if len(g.Log.Output) > 0 {
		// add flag for log output from provided build command
		flags = append(flags, fmt.Sprintf("--log-output=%s", g.Log.Output))
	}

	return flags
}
