// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"
	"reflect"
	"testing"
)

func TestMakisu_Push_Command(t *testing.T) {
	// setup types
	p := &Push{
		Path:           "/path/to/tar",
		Pushes:         []string{"FOO"},
		RegistryConfig: "{}",
		Replicas:       []string{"FOO"},
		Tag:            "latest",
	}

	// nolint // this functionality is not exploitable the way
	// the plugin accepts configuration
	want := exec.Command(
		_makisu,
		pushAction,
		fmt.Sprintf("--push \"%s\"", p.Pushes[0]),
		fmt.Sprintf("--registry-config %s", p.RegistryConfig),
		fmt.Sprintf("--replica \"%s\"", p.Replicas[0]),
		fmt.Sprintf("--tag %s", p.Tag),
		".",
	)

	got := p.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestMakisu_Push_Exec_Error(t *testing.T) {
	// setup types
	p := &Push{}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestMakisu_Push_Validate(t *testing.T) {
	// setup types
	p := &Push{
		Path: "/path/to/tar",
		Tag:  "latest",
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Push_Validate_NoContextPath(t *testing.T) {
	// setup types
	p := &Push{
		Tag: "latest",
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Push_Validate_NoTag(t *testing.T) {
	// setup types
	p := &Push{
		Path: "/path/to/tar",
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
