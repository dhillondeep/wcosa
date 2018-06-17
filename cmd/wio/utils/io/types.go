// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package io contains helper functions related to io
// This file contains all the types for IO that the user can use

package io

type IOHandler interface {
	GetRoot() (string, error)
	CopyFile(string, string, bool) error
	CopyMultipleFiles([]string, []string, []bool) error
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error
	ParseJson(string, interface{}) error
	ParseYml(string, interface{}) error
	WriteJson(string, interface{}) error
	WriteYml(string, interface{}) error
}

type NormalHandler byte
type AssetHandler byte
