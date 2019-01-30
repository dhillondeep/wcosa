package create

import (
    "os"
    "wio/internal/constants"
    "wio/internal/env"
    "wio/internal/types"
    "wio/pkg/util"
    "wio/pkg/util/sys"
    "wio/pkg/util/template"
)

func getIdeaFilePath(directory, fileName string) string {
    return sys.Path(directory, sys.IdeaFolder, fileName)
}

type dispatchIdeGenFunc func(directory string, target types.Target, config types.Config) (types.Target, error)

var dispatchIdeGenFunctions = map[string]dispatchIdeGenFunc{
    constants.Clion: dispatchIdeGenClion,
}

func (create Create) generateIdeFiles(ide string, directory string, config types.Config) (types.Target, error) {
    defTargetName := config.GetInfo().GetOptions().GetDefault()

    if util.IsEmptyString(defTargetName) {
        return nil, util.Error("Default Target must always be defined")
    }

    var defTarget types.Target
    var exists bool
    if defTarget, exists = config.GetTargets()[defTargetName]; exists {
        defTarget.SetName(defTargetName)
    } else {
        return nil, util.Error("Default Target %s is not defined in wio.yml file", defTargetName)
    }

    if val, exists := dispatchIdeGenFunctions[ide]; exists {
        return val(directory, defTarget, config)
    }

    return defTarget, nil
}

func dispatchIdeGenClion(directory string, target types.Target, config types.Config) (types.Target, error) {
    cmakeListsPath := sys.Path(directory, "CMakeLists.txt")
    miscPath := getIdeaFilePath(directory, "misc.xml")
    watcherTasksPath := getIdeaFilePath(directory,"watcherTasks.xml")
    workplacePath := getIdeaFilePath(directory,"workspace.xml")

    if err := os.MkdirAll(sys.Path(directory, sys.WioFolder, sys.IdeFolder, constants.Clion,
        target.GetName()), os.ModePerm); err != nil {
        return nil, err
    }

    if err := template.IOReplace(miscPath, map[string]string{
        "TARGET_SRC": target.GetSource(),
    }); err != nil {
        return nil, err
    }

    if err := template.IOReplace(watcherTasksPath, map[string]string{
        "WIO_PATH": env.GetWioPath(),
    }); err != nil {
        return nil, err
    }

    if err := template.IOReplace(workplacePath, map[string]string{
        "PROJECT_NAME": config.GetName(),
        "TARGET_NAME": target.GetName(),
        "TARGET_SRC": target.GetSource(),
    }); err != nil {
        return nil, err
    }

    return target, template.IOReplace(cmakeListsPath, map[string]string{
        "TARGET_NAME": target.GetName(),
        "WIO_PATH": env.GetWioPath(),
    })
}
