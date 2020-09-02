// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestMakisu_Config_Create(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		URL:      "index.docker.io",
		Username: "octocat",
	}

	_, err := c.Create()
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestMakisu_Config_Create_NoName(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		URL:      "index.docker.io",
		Username: "octocat",
	}

	_, err := c.Create()
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestMakisu_Config_Create_NoUsername(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		URL:      "index.docker.io",
		Username: "octocat",
	}

	_, err := c.Create()
	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestMakisu_Config_Create_NoPassword(t *testing.T) {
	// setup types
	c := &Config{
		Password: "superSecretPassword",
		URL:      "index.docker.io",
		Username: "octocat",
	}

	_, err := c.Create()
	if err != nil {
		t.Errorf("Create returned err: %v", err)
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
