## Wio

Wio is an open source Iot development environment. It let's user build c/c++ projects for various platforms in no time. At the moment only Atmel AVR platform is supported but, in future releases, and before this project is stable and out of beta, support for native (desktop c/c++) platform will be added.

### Contribute
If you are interested in working on wio, you can read [contribution document](https://github.com/wio/wio/blob/master/CONTRIBUTING.md) and add features/fixes.

Table of contents
=================

<!--ts-->
   * [Wio Install](#wio-install)
        * [Ubuntu](#ubuntu)
        * [Arch](#arch)
        * [MacOs](#macos)
        * [Windows](#windows)
        * [Github Releases](#github-releases)
   * [Create and Update](#create-and-update)
        * [App and Pkg](#app-and-pkg)
            * [Pkg](#pkg)
            * [App](#app)
        * [Update](#update)
   * [Build](#build)
   * [Upload](#upload)
   * [Devices](#devices)
   * [Package Manager](#package-manager)
        * [Publish](#publish)
        * [Install](#install)
<!--te-->

## Wio Install
Wio is available on all popular platforms. You can check a guide specific to your operating system below.

### Ubuntu
Best way to install wio on ubuntu is to use `NPM`. NPM is a node package manager and can be downloaded using:
```bash
sudo apt-get install npm
```
Then to install wio, you have to do:
```bash
npm install -g wio
```
Many times there are issues with permission on ubuntu and hence, following can be done:
```bash
sudo npm install -g wio --unsafe-perm
```
\
Since Wio is a development environment, you have to install compilers and build tools.
```bash
# This is for atmel AVR platform
sudo apt-get install gcc-avr avr-libc avrdude

# This is for building any wio project
sudo apt-get install cmake make
```

**Note:** All the framework files (SDK) are already included with wio and you do not have to worry about them.

### Arch
Best way to install wio on Arch is to use `NPM`. NPM is a node package manager and can be downloaded using:
```bash
sudo pacman -S nodejs
```
Then to install wio, you have to do:
```bash
npm install -g wio
```
Many times there are issues with permission on Arch and hence, following can be done:
```bash
sudo npm install -g wio --unsafe-perm
```
\
Since Wio is a development environment, you have to install compilers and build tools.
```bash
# This is for atmel AVR platform
sudo pacman -S avr-gcc avr-libc avrdude

# This is for building any wio project
sudo pacman -S cmake make
```

**Note:** All the framework files (SDK) are already included with wio and you do not have to worry about them.

### MacOs
There are two good ways to install wio on MacOS. You can use `NPM` or `Homebrew`.

**NPM**

NPM is a node package manager and can be downloaded using:
```bash
brew install node
```
Then to install wio, you have to do:
```bash
npm install -g wio
```
If there are issues with permissions, following can be done:
```bash
sudo npm install -g wio --unsafe-perm
```

**Homebrew**

Homebrew is a package manager for macOS. In order to install wio, you have to do:
```bash
brew tap dhillondeep/wio
brew install dhillondeep/wio/wio
```
If you are planning on using [wio package manager](#package-manager), you will need `npm` installed. To install `npm`, you have to do:
```bash
brew install node
```
\
Since Wio is a development environment, you have to install compilers and build tools.
```bash
# This is for atmel AVR platform
xcode-select --install
brew tap osx-cross/avr
brew install avr-gcc
brew install avrdude

# This is for building any wio project
sudo pacman -S cmake make
```

**Note:** All the framework files (SDK) are already included with wio and you do not have to worry about them.

### Windows
There are multiple ways to install wio. Instructions for each are mentioned below.

**Scoop**

Scoop is a windows package manager. It downloads binaries like how linux package managers do. You can use scoop if you have powershell 13 or above. To install Scoop, type the following in powershell:
```bash
iex (new-object net.webclient).downloadstring('https://get.scoop.sh')
```
If you get an error you might need to change the execution policy (i.e. enable Powershell):
```bash
set-executionpolicy -s cu unrestricted
```
After installing Scoop, you can install wio:
```bash
scoop bucket add wio https://github.com/dhillondeep/wio-bucket.git
scoop install wio
```

**NPM**

NPM is a node package manager and can be downloaded through [nodejs website](https://nodejs.org/en/download/)

Then to install wio, you have to do:
```bash
npm install -g wio
```

\
Since Wio is a development environment, you have to install compilers and build tools.

**Cmake and Make**

You can use Scoop to download `cmake` and `make`:
```bash
scoop install cmake make
```
Or you can download CMake from [CMake website](https://cmake.org/download/) and Make from [GNU32 Website](http://gnuwin32.sourceforge.net/packages/make.htm). Make sure to add these tools to your Environment PATH

**Atmel AVR toolchain**

* Download [Core Arduino](https://www.arduino.cc/en/Main/Software) and install it. This will contains all the compilers and tools needed for building and uploading the binaries
* Go to Arduino install directory, ex: `C:\Program Files (x86)\Arduino` and from there go to `hardware\tools\avr\bin`. Grab the whole path and add that to your environment PATH

**Note:** All the framework files (SDK) are already included with wio and you do not have to worry about them.

## Github Releases
Alternatively, all the versions are pushed to Github Releases and can be downloaded from [Github Releases site](https://github.com/wio/wio/releases).

## Create and Update
To learn about create or update command
```bash
wio create -h
wio update -h
```

`Wio` enforces a strict project structure so that development is organized and easier to build. This project structure will be created by `create` command and update by `update` commands. To create a wio project:

```bash
# for project of type app
wio create app <platform> [directory] [board]
```
```bash
# for project of type pkg
wio create pkg <platform> [directory] [board]
```

If directory is not provided, current directory is used. So far the only supported platform is `avr`

### App and Pkg
As you can see in the example above, there are two types of projects you can create. This is a feature of wio that provides you flexibility when working with iot development. Our goal is to make development modular by creating `packages` that other wio projects can import. These packages are then imported by project of type `app` (executable projects). This way code is easily transferable. Below is in depth about what you can do with project of type `pkg` and of type `app`.

#### Pkg
Typical Structure for `pkg` project is as follows:
```
.
├── include
│   └── someFile.h
├── src
│   └── someFile.cpp
├── test
│   └── pkgTest.cpp
└── wio.yml
```

Notes about the structure:
* Always put public headers in `include` directory.
* Anything private (headers or not) goes in `src` folder.
* Since this is a package, `test` directory contains code that uses this package to create an executable.
    * Package is automatically linked and headers are included
    * You can test you application this way.
    * This directory is not pushed when this package is published.

`wio.yml` allows you to provide meta data for the package and all sorts of configurations. A typical content looks like this:

```yml
pkg:
  ide: none

  meta:
    name: example
    description: A wio avr pkg using cosa framework
    version: 0.0.1
    keywords:
    - avr
    - c
    - c++
    - wio
    - cosa
    license: MIT

  config:
    minimum_wio_version: 0.3.2
    supported_platforms:
    - avr
    supported_frameworks:
    - cosa
    supported_boards:
    - uno

  compile_options:
    header_only: false
    platform: avr
  flags:
    allow_only_global_flags: false
    allow_only_required_flags: false
    global_flags: []
    visibility: PRIVATE
  definitions:
    allow_only_global_definitions: false
    allow_only_required_definitions: false
    global_definitions: []
    visibility: PRIVATE

targets:
  default: tests
  create:
    tests:
      src: tests
      framework: cosa
      board: uno

dependencies:
  package1:
    version: 0.0.1
```

#### App
Typical Structure for `app` project is as follows:

```
.
├── src
│   ├── someFile.cpp
│   └── someFile.h
└── wio.yml
```

Notes about the structure:
* All the code goes in `src` directory. No other directory is used.
* In future versions we will add support for testing and then tests will be in `tests` directory.

`wio.yml` is used to configure executable development and dependency management for the project. A typical content looks like this:

```yml
app:
  name: example
  ide: none

  config:
    minimum_wio_version: 0.3.2
    supported_platforms:
    - avr
    supported_frameworks:
    - cosa
    supported_boards:
    - uno

  compile_options:
    platform: avr

targets:
  default: main
  create:
    main:
      src: src
      framework: cosa
      board: uno

dependencies:
  package1:
    version: 0.0.1
```

### Update
`update` is a compimentary command to `create` and it allows for updating anything "broken" inside your project. For example, if you rename the directory, you need to run update command to make sure the name of the project is updated. Another example is that if you change `framework`, `platform`, you need to run update. This command is also used if `wio.yml` file is broken and a new one needs to be created. To update a wio project:

```bash
# for project of type app
wio update [directory]

# for project of type pkg
wio update [directory]

# for only creating wio.yml file
wio update --only-config [directory]
```

If directory is not provided, current directory is used. So far the only supported platform is `avr`

## Build
To learn about build command
```bash
wio run -h
```

Building a wio project is very easy and it comes with many configurations you can uses. Here are examples:

```bash
# build default target
wio run [directory]
```
```bash
# build specific target
wio run [directory] -target target2
```
```bash
# clean build files before building
wio run [directory] --clean
```

If directory is not provided, current directory is used. So far the only supported platform is `avr`

All of these commands above will build the project. Wio takes care of all the build files, dependency trees, etc. You just need to put files in proper project structure and all will be good. Wio uses `CMake` and `Make` to build the project. These build files are generated when `run` command is executed and hence the latest changes are picked up from `wio.yml` file. All the build files and internal helpers are stored in `.wio` folder. This folder should never be touched and modified since it can cause unexpected behaviour.

## Upload
To learn about upload command
```bash
wio run -h
```

`upload` command is used by platforms that need the executable to be uploaded to some hardware. Uploading is simple:


```bash
# normal upload (default target)
wio run --upload
```
```bash
# build with upload (different target)
wio run --upload --target target2
```
```bash
# build and upload to specified port
wio run --upload --port /dev/ttySu
```

## Devices
To learn about devices command
```bash
wio devices -h
```

This gives user information about the ports and devices connected to them. It also allows user to open a serial monitor. Examples are below:

```bash
# list open devices
wio devices list
```
```bash
# opens serial monitor (automatic port detection for arduino official boards)
wio devices monitor
```
```bash
# opens serial monitor for port specifed
wio devices monitor --port /dev/ttySu
```

## Package Manager
A goal of Wio has always been to make iot development easier and modular. No more cloning the repository and linking the code. This is achieved by wio package manager called. Wio is using `npm` package manager as a backend and a lot of wrapper code has been written around it to make it a c/c++ package manager. Let's talk about command

### Publish
To learn about package manager's publish command
```bash
wio publish -h
```

`publish` command is used when your `pkg` project is ready and now you want to publish it for other people to use. You can fill in all the meta information and other development information in `wio.yml` file. Since wio is using npm, user has to login in order to publish a package.

```bash
# publish pkg from current directory
wio publish [directory]
```

```bash
# will create an account or log you in (ignore if already logged in)
npm addusr
```

### Install
To learn about package manager's install command
```bash
wio install -h
```

`install` command is used to pull in all the dependencies for the project. When this command is executed, wio goes through all the listed dependencies in `wio.yml` file and pulls all of them recursively. In the build process, these dependencies are linked to compile the project. Example of `install` command

```bash
# pulls all the dependencies for the current directory
wio install [package name]
```
```bash
# pulls dependencies and add them to wio.yml file
wio install package1 package2 --save
```
```bash
# pulls all the dependencies and clean the old ones
wio install --clean
```

If package name is provided, only that package is pulled regardless if it is in `wio.yml` file. This command also acts as an update command as well. You can modify the versions of dependencies in `wio.yml` file and then run this command. This will automatically upgrade or downgrade the dependency and update the dependency tree.
