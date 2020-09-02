// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestMakisu_Plugin_Exec(t *testing.T) {
	// TODO Write test
}

func TestMakisu_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Password: "superSecretPassword",
			URL:      "index.docker.io",
			Username: "octocat",
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
