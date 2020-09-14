// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const pushAction = "push"

type (
	// Push represents the plugin configuration for push information.
	//
	// Makisu documents their command usage:
	// https://githup.com/uber/makisu/blob/master/docs/COMMAND.md
	// nolint
	Push struct {
		// enables setting the location of the tar file
		Path string
		// enables setting registries to push an image to
		Pushes []string
		// enables setting the credentials for authentication
		RegistryConfig string
		// enables pushing images with alternative full names
		Replicas []string
		// enables setting the tag for an image
		Tag string
	}
)

// pushFlags represents for config settings on the cli.
// nolint
var pushFlags = []cli.Flag{
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_PATH"},
		FilePath: string("/vela/parameters/makisu/push/path,/vela/secrets/makisu/push/path"),
		Name:     "push.context",
		Usage:    "enables setting the location of the tar file",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_PUSHES"},
		FilePath: string("/vela/parameters/makisu/push/context,/vela/secrets/makisu/push/context"),
		Name:     "push.pushes",
		Usage:    "enables setting registries to push an image to",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_REGISTRY_CONFIG"},
		FilePath: string("/vela/parameters/makisu/push/registry_config,/vela/secrets/makisu/push/registry_config"),
		Name:     "push.regstry-config",
		Usage:    "enables setting the credentials for authentication",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"PARAMETER_CONTEXT"},
		FilePath: string("/vela/parameters/makisu/push/context,/vela/secrets/makisu/push/context"),
		Name:     "push.replicas",
		Usage:    "enables pushing images with alternative full names",
	},
	&cli.StringFlag{
		EnvVars:  []string{"PARAMETER_TAG"},
		FilePath: string("/vela/parameters/makisu/push/registry_config,/vela/secrets/makisu/push/registry_config"),
		Name:     "push.tag",
		Usage:    "enables setting the tag for an image",
	},
}

// Command formats and outputs the Push command from
// the provided configuration to push a Docker image.
// nolint
func (p *Push) Command() *exec.Cmd {
	logrus.Trace("creating makisu push command from plugin configuration")

	// variable to store flags for command
	var flags []string

	// check if Pushes is provided
	if len(p.Pushes) > 0 {
		var args string
		for _, arg := range p.Pushes {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for Pushes from provided push command
		flags = append(flags, fmt.Sprintf("--push \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// check if RegistryConfig is provided
	if len(p.RegistryConfig) > 0 {
		// add flag for RegistryConfig from provided push command
		flags = append(flags, fmt.Sprintf("--registry-config %s", p.RegistryConfig))
	}

	// check if Replicas is provided
	if len(p.Replicas) > 0 {
		var args string
		for _, arg := range p.Replicas {
			args += fmt.Sprintf(" %s", arg)
		}
		// add flag for Replicas from provided push command
		flags = append(flags, fmt.Sprintf("--replica \"%s\"", strings.TrimPrefix(args, " ")))
	}

	// check if Tag is provided
	if len(p.Tag) > 0 {
		// add flag for Tag from provided push command
		flags = append(flags, fmt.Sprintf("--tag %s", p.Tag))
	}

	// add the required directory param
	flags = append(flags, p.Path)

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	return exec.Command(_makisu, append([]string{pushAction}, flags...)...)
}

// Exec formats and runs the commands for pushing a Docker image.
func (p *Push) Exec() error {
	logrus.Trace("running push with provided configuration")

	// create the push command for the file
	cmd := p.Command()

	// run the push command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Push is properly configured.
func (p *Push) Validate() error {
	logrus.Trace("validating push plugin configuration")

	// verify tag are provided
	if len(p.Path) == 0 {
		return fmt.Errorf("no push path provided")
	}

	// verify tag are provided
	if len(p.Tag) == 0 {
		return fmt.Errorf("no push tag provided")
	}

	return nil
}
