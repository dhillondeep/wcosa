package create

import "github.com/urfave/cli"

type Create struct {
    Context *cli.Context
    Update  bool
    error   error
}

type createInfo struct {
    Directory string
    Type      string
    Name      string

    Platform  string
    Framework string
    Board     string

    ConfigOnly bool
    HeaderOnly bool
}

// get context for the command
func (create Create) GetContext() *cli.Context {
    return create.Context
}

// Executes the create command
func (create Create) Execute() {
    directory := performDirectoryCheck(create.Context)

    if create.Update {
        // this checks if wio.yml file exists for it to update
        performWioExistsCheck(directory)
        // this checks if project is valid state to be updated
        performPreUpdateCheck(directory, &create)
        create.handleUpdate(directory)
    } else {
        // this checks if directory is empty before create can be triggered
        performPreCreateCheck(directory, create.Context.Bool("only-config"))
        create.createPackageProject(directory)
    }
}
