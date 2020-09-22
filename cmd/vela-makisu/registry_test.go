// Copyright (c) 2020 Target Brands, Inr. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestMakisu_Registry_Write(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		DryRun:   false,
		Password: "superSecretPassword",
		Addr:     "index.docker.io",
		Username: "octocat",
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestMakisu_Registry_Validate(t *testing.T) {
	// setup types
	r := &Registry{
		Password: "superSecretPassword",
		Addr:     "index.docker.io",
		Username: "octocat",
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Registry_Validate_NoPassword(t *testing.T) {
	// setup types
	r := &Registry{
		Addr:     "index.docker.io",
		Username: "octocat",
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Registry_Validate_NoURL(t *testing.T) {
	// setup types
	r := &Registry{
		Password: "superSecretPassword",
		Username: "octocat",
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Registry_Validate_NoUsername(t *testing.T) {
	// setup types
	r := &Registry{
		Password: "superSecretPassword",
		Addr:     "index.docker.io",
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
