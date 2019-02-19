# Development Setup

## Building Wio with Your Changes

Wio uses [mage](https://github.com/magefile/mage) to sync vendor dependencies, build Wio, and other things. You must run mage from the Wio directory.

Install mage
```bash
go get github.com/magefile/mage
```

To build wio
```bash
mage build
```

Wio binary is created inside the `bin` fodler of your root directory. You can source `wenv` for unix based systems and `env.bat` or `env.ps1` for windows.
This will add the binary to the path.

To build wio
```bash
mage build
```

To clean build files
```bash
mage clean
```

To install dependencies
```bash
mage install
```

Currently the script only runs tests on `bash` based systems. To run the tests:
```bash
wmake test
```

To list all available commands along with descriptions:
```bash
mage -l
```


### Installing Go

Wio development should be done on Go version 1.11

