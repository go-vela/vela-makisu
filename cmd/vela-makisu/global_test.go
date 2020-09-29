// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMakisu_Global_Flag(t *testing.T) {
	// setup types
	g := &Global{
		CPU: &CPU{
			Profile: true,
		},
		Log: &Log{
			Fmt:    "console",
			Level:  "debug",
			Output: "stdout",
		},
	}

	want := []string{
		"--cpu-profile",
		fmt.Sprintf("--log-fmt=%s", g.Log.Fmt),
		fmt.Sprintf("--log-level=%s", g.Log.Level),
		fmt.Sprintf("--log-output=%s", g.Log.Output),
	}

	got := g.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flag is %v, want %v", got, want)
	}
}
