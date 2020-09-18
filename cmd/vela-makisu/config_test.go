// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestMakisu_Config_Write(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	c := &Config{
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "index.docker.io",
		Username: "octocat",
	}

	err := c.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestMakisu_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		URL:      "index.docker.io",
		Username: "octocat",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Config_Validate_NoPassword(t *testing.T) {
	// setup types
	c := &Config{
		URL:      "index.docker.io",
		Username: "octocat",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Config_Validate_NoURL(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		Username: "octocat",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Config_Validate_NoUsername(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		URL:      "index.docker.io",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
