// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
	"time"
)

func TestMakisu_Plugin_Exec(t *testing.T) {
	// TODO Write test
}

func TestMakisu_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Registry: &Registry{
			Password: "superSecretPassword",
			Addr:     "index.docker.io",
			Username: "octocat",
		},
		//nolint
		Build: &Build{
			BuildArgs:   []string{"FOO"},
			Context:     ".",
			Commit:      "b0bb040e6a6d71ddf98684349c42d36fa6c539ad",
			Compression: "default",
			DenyList:    []string{"FOO"},
			Docker: &Docker{
				Host:    "unix:///var/run/docker.sock",
				Scheme:  "http",
				Version: "1.21",
			},
			Destination: "/path/to/dest",
			File:        "Dockerfile",
			HTTPCache: &HTTPCache{
				Addr:    "http://localhost",
				Headers: []string{"Content-type: Application/json"},
			},
			Load:          true,
			LocalCacheTTL: 1 * time.Minute,
			ModifyFS:      true,
			PreserveRoot:  true,
			Pushes:        []string{"FOO"},
			RedisCache: &RedisCache{
				Addr:     "http://localhost",
				Password: "superSecret123",
				TTL:      1 * time.Minute,
			},
			RegistryConfig: "{}",
			Replicas:       []string{"FOO"},
			Storage:        "foo",
			Tag:            "latest",
			Target:         "dev",
		},
		Push: &Push{
			Path:           ".",
			Pushes:         []string{"FOO"},
			RegistryConfig: "{}",
			Replicas:       []string{"FOO"},
			Tag:            "latest",
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
