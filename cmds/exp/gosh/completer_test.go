// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/u-root/prompt"
)

func TestCompleterFunc(t *testing.T) {
	for _, tt := range []struct {
		name      string
		inputText string
		env       string
		resultSet []string
	}{
		{
			name:      "no text",
			inputText: "",
			resultSet: []string{},
		},
		{
			name:      "echo",
			inputText: "ec",
			env:       "/bin",
			resultSet: []string{"echo"},
		},
		{
			name:      "pipe",
			inputText: "echo test | ec",
			resultSet: []string{},
		},
		{
			name:      "files",
			inputText: "./",
			resultSet: []string{"completer_test.go", "completer.go"},
		},
		{
			name:      "wrong path",
			inputText: "ec",
			env:       "/bogus",
			resultSet: []string{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b := prompt.NewBuffer()
			b.InsertText(tt.inputText, false, true)

			origPath := os.Getenv("PATH")
			if err := os.Setenv("PATH", tt.env); err != nil {
				t.Errorf("Failed setting environment: %v", err)
			}
			suggestions := completerFunc(*b.Document())
			if !contentsEqual(tt.resultSet, suggestions) {
				t.Errorf("want: %v got: %v", tt.resultSet, suggestions)
			}
			if err := os.Setenv("PATH", origPath); err != nil {
				t.Errorf("Failed resetting environment: %v", err)
			}
		})
	}
}

func contentsEqual(want []string, got []prompt.Suggest) bool {
	for _, entry := range want {
		found := false
		for _, suggestion := range got {
			if entry == suggestion.Text {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
