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

func TestImg_Push_Command(t *testing.T) {
	// setup types
	b := &Push{
		ContextPath:    ".",
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
		fmt.Sprintf("--push \"%s\"", b.Pushes[0]),
		fmt.Sprintf("--registry-config %s", b.RegistryConfig),
		fmt.Sprintf("--replica \"%s\"", b.Replicas[0]),
		fmt.Sprintf("--tag %s", b.Tag),
		".",
	)

	got := b.Command()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestMakisu_Push_Exec_Error(t *testing.T) {
	// setup types
	b := &Push{}

	err := b.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestMakisu_Push_Validate(t *testing.T) {
	// setup types
	b := &Push{
		ContextPath: ".",
		Tag:         "latest",
	}

	err := b.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestMakisu_Push_Validate_NoContextPath(t *testing.T) {
	// setup types
	b := &Push{
		Tag: "latest",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestMakisu_Push_Validate_NoTag(t *testing.T) {
	// setup types
	b := &Push{
		ContextPath: ".",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
