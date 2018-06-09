// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Types for all the things being used in this package
package create

type Data struct {
    Id       string
    Src      string
    Des      string
    Override bool
}

type Paths struct {
    Paths []Data
}

/// This structure wraps all the important features needed for a create and update command
type PacketCreate struct {
    ProjType          string
    Update            bool
    Directory         string
    Name              string
    Board             string
    Framework         string
    Platform          string
    Ide               string
    Tests             bool
    CreateDemo        bool
    CreateExtras      bool
    HeaderOnlyFlagSet bool
    HeaderOnly        bool
}

// #################################### Type for project structure json ##################################
type StructureFilesData struct {
    Constrains []string
    From       string
    To         string
    Override   bool
    Update     bool
}

type StructurePathData struct {
    Constrains []string
    Entry      string
    Files      []StructureFilesData
}

type StructureTypeData struct {
    Paths []StructurePathData
}

type StructureConfigData struct {
    App StructureTypeData
    Pkg StructureTypeData
}
