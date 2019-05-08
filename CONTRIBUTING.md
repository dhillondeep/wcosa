# Development Setup

Wio is developed using go and version being used is 1.12.  Wio uses [task](https://github.com/go-task/task) 
to build the project. You will need to have this installed. Installing instructions for this tool can be 
found on https://taskfile.dev/#/installation

#### Building wio
```bash
task templateGen
task build
```

#### Testing
Currently there are only unit tests that can be run
```bash
task templateGen
task unitTest
```

## Adding features and fixes
In order to create a feature, you will need to create new branch explaining the feature. The code inside the branch
should have full test coverage. Then PR can be made against `develop-1.0.0` branch. After all the tests have been
passed, code will be merged.
