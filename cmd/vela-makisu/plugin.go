// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
	// Used for translating the raw docker configuration
	Global *Global
	// enables setting configuration for the global flags
	GlobalRaw string
	// registry arguments loaded for the plugin
	Registry *Registry
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
	err = p.Registry.Write()
	if err != nil {
		return err
	}

	// get any global flags that may have been set
	globalFlags := p.Global.Flags()

	// set any configuration for global flags
	p.Build.GlobalFlags = globalFlags

	// set required configuration for registry config
	p.Build.RegistryConfig = configPath

	// execute build action
	return p.Build.Exec()
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (p *Plugin) Unmarshal() error {
	logrus.Trace("unmarshaling global flags")

	// allocate configuration to structs
	p.Global = &Global{}

	// check if any global flags were passed
	if len(p.GlobalRaw) > 0 {
		// cast raw global flags into bytes
		globalFlags := []byte(p.GlobalRaw)

		// serialize raw global flags into expected Global type
		err := json.Unmarshal(globalFlags, &p.Global)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate verifies the Plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate config configuration
	err := p.Registry.Validate()
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
