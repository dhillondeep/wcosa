// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Helper for create and update command to update targets
package create

import (
    "strings"
    "wio/cmd/wio/config"
    "wio/cmd/wio/types"
)

// Updates AVR App Targets to make sure there is atleast one valid target
func updateAVRAppTargets(targets *types.AppTargets, directory string) {
    //////////////////////////////////////////// Targets //////////////////////////////////////////////////
    if strings.Trim(targets.DefaultTarget, " ") != "" {
        // check if default target does not exist
        if _, exists := targets.Targets[targets.DefaultTarget]; !exists {
            // create a default target
            targets.Targets[targets.DefaultTarget] = &types.AppTarget{
                Src:       "src",
                Framework: config.ProjectDefaults.Framework,
                Board:     config.ProjectDefaults.AVRBoard,
                Flags: types.AppTargetFlags{
                    GlobalFlags: []string{},
                    TargetFlags: []string{},
                },
            }
        }
    } else {
        // create a default target or make one the default
        if len(targets.Targets) <= 0 {
            targets.DefaultTarget = config.ProjectDefaults.AppTargetName
            // create a default target
            targets.Targets = map[string]*types.AppTarget{
                config.ProjectDefaults.AppTargetName: {
                    Src:       "src",
                    Framework: config.ProjectDefaults.Framework,
                    Board:     config.ProjectDefaults.AVRBoard,
                    Flags: types.AppTargetFlags{
                        GlobalFlags: []string{},
                        TargetFlags: []string{},
                    },
                },
            }
        } else {
            for targetName := range targets.Targets {
                targets.DefaultTarget = targetName
                break
            }
        }
    }
}
