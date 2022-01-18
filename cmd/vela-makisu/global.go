// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type (
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
	// configFlags represents the global flags that can be set on the makisu commands.
	// nolint
	globalFlags = []cli.Flag{
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_GLOBAL_FLAGS"},
			FilePath: string("/vela/parameters/makisu/global/flags,/vela/secrets/makisu/global/flags"),
			Name:     "global.flags",
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

// Flags formats and outputs the flags for
// configuring a Docker.
func (g *Global) Flags() []string {
	logrus.Trace("creating global flags for command")

	// variable to store flags for command
	var flags []string

	// check if cpu profile is provided
	if g.CPU.Profile {
		// add flag for cpu profile from provided build command
		flags = append(flags, "--cpu-profile")
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
