package create

import "github.com/urfave/cli"

type Create struct {
    Context *cli.Context
    Update  bool
    Type    string
}

type createInfo struct {
    context *cli.Context

    directory   string
    projectType string
    name        string

    platform  string
    framework string
    board     string

    configOnly bool
    headerOnly bool
    updateOnly bool
}

// get context for the command
func (create Create) GetContext() *cli.Context {
    return create.Context
}

// Executes the create command
func (create Create) Execute() error {
    directory, err := performDirectoryCheck(create.Context)
    if err != nil {
        return err
    }

    if create.Update {
        // this checks if wio.yml file exists for it to update
        if err := performWioExistsCheck(directory); err != nil {
            return err
        }
        // this checks if project is valid state to be updated
        if err := performPreUpdateCheck(directory, &create); err != nil {
            return err
        }
        return create.handleUpdate(directory)
    } else {
        // this checks if directory is empty before create can be triggered
        onlyConfig := create.Context.Bool("only-config")
        if err := performPreCreateCheck(directory, onlyConfig); err != nil {
            return err
        }
        return create.createProject(directory)
    }
    return nil
}
