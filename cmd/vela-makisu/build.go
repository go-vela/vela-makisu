// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const buildAction = "build"

// Build represents the plugin configuration for build information.
type Build struct {
}

// buildFlags represents for config settings on the cli.
var buildFlags = []cli.Flag{}

// Command formats and outputs the Build command from
// the provided configuration to build a Docker image.
func (b *Build) Command() *exec.Cmd {
	logrus.Trace("creating makisu build command from plugin configuration")

	return nil
}

// Exec formats and runs the commands for building a Docker image.
func (b *Build) Exec() error {
	logrus.Trace("running build with provided configuration")

	// create the build command for the file
	cmd := b.Command()

	// run the build command for the file
	err := execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	return nil
}
