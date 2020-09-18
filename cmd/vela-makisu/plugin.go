// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
	// config arguments loaded for the plugin
	Config *Config
	// push arguments loaded for the plugin
	Push *Push
}

// Exec formats and runs the commands for building and publishing a Docker image.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	// output makisu version for troubleshooting
	err := execCmd(versionCmd())
	if err != nil {
		return err
	}

	// create config configuration for authentication to a registry
	err = p.Config.Write()
	if err != nil {
		return err
	}

	// get any global flags that may have been set
	globalFlags := p.Config.Global.Flags()

	// set any configuration for global flags
	p.Build.GlobalFlags = globalFlags
	p.Push.GlobalFlags = globalFlags

	// set required configuration for registry config
	p.Build.RegistryConfig = configPath
	p.Push.RegistryConfig = configPath

	// execute build action
	path, err := p.Build.Exec()
	if err != nil {
		return err
	}

	// set the location to the built image
	p.Push.Path = path

	// execute push action if not in dry run mode
	if !p.Config.DryRun {
		// validate push configuration
		err = p.Push.Validate()
		if err != nil {
			return err
		}

		// execute push action
		err = p.Push.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate verifies the Plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// when user adds global flag configuration
	err := p.Config.Unmarshal()
	if err != nil {
		return err
	}

	// validate config configuration
	err = p.Config.Validate()
	if err != nil {
		return err
	}

	// when user adds configuration additional options
	// for: docker, http, redis
	err = p.Build.Unmarshal()
	if err != nil {
		return err
	}

	// validate build configuration
	err = p.Build.Validate()
	if err != nil {
		return err
	}

	return nil
}
