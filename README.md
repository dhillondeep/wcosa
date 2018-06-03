[![Build Status](https://travis-ci.org/waterloop/wcosa.svg?branch=master)](https://travis-ci.org/waterloop/wcosa)

## Wio
Wio is a AVR development environment which let's you create, build, test and upload AVR programs. It current only
supports `Cosa` SDK but, the support for `Arduino` SDK will be added in near future. Similarly, it only supports
`AVR` development but, `ATMEL` boards support will be added in near future as well. Read the documentation below for features and uses.

Table of contents
=================

<!--ts-->
   * [Install](#install)
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
   * [Clean](#clean)
   * [Run](#run)
   * [Package Manager](#package-manager)
        * [Publish](#publish)
        * [Get](#get)
        * [Vendor](#vendor)
   * [Toolchain](#toolchain)
<!--te-->

## Install
Wio is available on all popular platforms. You can check a guide specific to your operating system below.

### Ubuntu
Best way to install wio on ubuntu is to use `NPM`. NPM is a node package manager and can be downloaded using:
```bash
sudo apt-get install npm
```
Then to install wio, you have to do:
```bash
sudo npm install -g wio --unsafe-perm
```
\
Since Wio is an embedded development environment, you have to install compilers and build tools. So far Wio only supports `AVR` development and hence `AVR` toolchain is needed. In order to install that, you can do:
* Download [Arduino](https://www.arduino.cc/en/Main/Software) or:
```bash
sudo apt-get install gcc-avr avr-libc avrdude
```
* CMake and Make
```bash
sudo apt-get install cmake make
```
All the framework files (SDK) are already included with wio and you do not have to worry about them.

### Arch
Best way to install wio on Arch is to use `NPM`. NPM is a node package manager and can be downloaded using:
```bash
sudo pacman -S nodejs
```
Then to install wio, you have to do:
```bash
sudo npm install -g wio --unsafe-perm
```
\
Since Wio is an embedded development environment, you have to install compilers and build tools. So far Wio only supports `AVR` development and hence `AVR` toolchain is needed. In order to install that, you can do:
* Download [Arduino](https://www.arduino.cc/en/Main/Software) or:
```bash
sudo pacman -S avr-gcc avr-libc avrdude
```
* CMake and Make
```bash
sudo pacman -S cmake make
```
All the framework files (SDK) are already included with wio and you do not have to worry about them.

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
Since Wio is an embedded development environment, you have to install compilers and build tools. So far Wio only supports `AVR` development and hence `AVR` toolchain is needed. In order to install that, you can do:
* Download [Arduino](https://www.arduino.cc/en/Main/Software) or:
```bash
# Keep in mind that building avr-gcc may take some time
xcode-select --install
brew tap osx-cross/avr
brew install avr-gcc
brew install avrdude
```
* CMake and Make
```bash
brew install cmake make
```
All the framework files (SDK) are already included with wio and you do not have to worry about them.

### Windows
There are multiple ways to install wio. Instructions for each are mentioned below.

### Scoop
Scoop is a windows package manager. It downloads binaries like how linux package managers do. You can use scoop if you have powershell 13 or above. To install Scoop, type the following in powershell:
```bash
iex (new-object net.webclient).downloadstring('https://get.scoop.sh')
```
**Note:** if you get an error you might need to change the execution policy (i.e. enable Powershell):
```bash
set-executionpolicy -s cu unrestricted
```
After installing Scoop, you can install wio:
```bash
scoop bucket add wio https://github.com/dhillondeep/wio-bucket.git
scoop install wio
```
\
Since Wio is an embedded development environment, you have to install compilers and build tools. You can use Scoop to download `cmake` and `make`:
```bash
scoop install cmake make
```
In order for [wio package manager](#package-manager) to work, you need to install `npm`. To do that:
```bash
scoop install nodejs
```
We still need toolchain for AVR compilers. For this,
* Download [Core Arduino](https://www.arduino.cc/en/Main/Software) and install it. This will contains all the compilers and tools needed for building and uploading the binaries
* Go to Arduino install directory, ex: `C:\Program Files (x86)\Arduino` and from there go to `hardware\tools\avr\bin`. Grab the whole path and add that to your environment PATH
* You will have to restart your shell for changes to take effect

All the framework files (SDK) are already included with wio and you do not have to worry about them.

### NPM
NPM is a node package manager and can be downloaded through [nodejs website](https://nodejs.org/en/download/)

Then to install wio, you have to do:
```bash
npm install -g wio
```
\
Since Wio is an embedded development environment, you have to install compilers and build tools. So far Wio only supports `AVR` development and hence `AVR` toolchain is needed. In order to install that, you follow the steps below:
* Download [Core Arduino](https://www.arduino.cc/en/Main/Software) and install it. This will contains all the compilers and tools needed for building and uploading the binaries
  * Go to Arduino install directory, ex: `C:\Program Files (x86)\Arduino` and from there go to `hardware\tools\avr\bin`. Grab the whole path and add that to your environment PATH
* Wio uses `CMake` for build files generation. Install it from [CMake website](https://cmake.org/download/)
  * Make sure to add `CMake` to your Environment Variables PATH
* Wio uses `Make` for building the project. There are multiple ways to get this:
  * If you have MinGW, you can add the `bin` folder to Environment PATH and this give you mingw32-make.exe
  * If you want a standalone `Make`, install it from [GNU32 Website](http://gnuwin32.sourceforge.net/packages/make.htm)
    * Make sure to add this to your Environment PATH

All the framework files (SDK) are already included with wio and you do not have to worry about them.

## Github Releases
Alternatively, all the versions are pushed to Github Releases and can be downloaded from [Github Releases site](https://github.com/dhillondeep/wio/releases).

## Create and Update
`Wio` enforces a strict project structure so that development is organized and easier to build. This project structure will be created by `create` and `update` commands provided. To create a wio project:

```bash
# for project of type app
wio create app <directory> <board>
```
```bash
# for project of type pkg
wio create pkg <directory> <board>
```

### App and Pkg
As you can see in the example above, there are two types of projects you can create. This is a feature of wio that provides you flexibility when working with embdedded development. Our goal is to make development modular by creating `packages` that other wio projects can import. These packages are then imported by project of type `app` (executable projects). This way code is easily transferable. Below is in depth about what you can do with project of type `pkg` and of type `app`.

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
    * You can test you application this way.
    * This directory is not pushed when this package is published.

`wio.yml` allows you to provide meta data for the package and all sorts of configurations. A typical content looks like this:

```yml
lib:
  # meta information
  name:
  description:
  url:
  version:
  author:
  contributors: []
  license:
  keywords: []

  # development flags
  platform:
  framework: []
  board: []
  compile_flags: []
  ide:

# targets to build for testing
targets:
  default: test
  created:
    test:
      board:
      compile_flags: []

# dependencies
dependencies:
  depName:
    version:
    vendor:
    compile_flags: []
```

As you can see above, there are three parts to this configuration file: 
* The information from tags from `meta information` is used when this package is published. Make sure you incerement the version everytime you push this package.
* The information from tags from `development flags` is used to see see how many platforms, frameworks and boards it supports. The `compile_flags` will be used when this package will be added as a dependency.  
    * `ide` is not supported at the moment. In future releases support for clion and VS Code will be added
    * `platform` is the chipset for the boards you are using. We only support `AVR` at the moment.
    * `framework` is the sdk you use to program. We only support `COSA` sdk but support for `Arduino` sdk is comming soon.
* The information from tags from `targets` is used to create targets that can be built and uploaded. By default a `test` target is created which can be used to test you package. You can create different targets with different `compile_flags` and `boards`.
* Dependencies are explained more below.

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
  name:
  ide:
  platform:
  framework:

targets:
  default: main
  created:
    main:
      board:
      compile_flags: []

```

Notes about the the configuration file:
* `ide` is not supported at the moment. In future releases support for clion and VS Code will be added
* `platform` is the chipset for the boards you are using. We only support `AVR` at the moment.
* `framework` is the sdk you use to program. We only support `COSA` sdk but support for `Arduino` sdk is comming soon.
* A default target main is created for you but you can create as many targets with different `boards` and `compile_flags`.

### Update
`update` is a compimentary command to `create` and it allows for updating anything "broken" inside your project. For example, if you rename the directory, you need to run update command to make sure the name of the project is updated. Another example is that if you change `framework`, `platform`, you need to run update. This command is also used if `wio.yml` file is broken and a new one needs to be created. To update a wio project:

```bash
# for project of type app
wio update app <directory>

# for project of type pkg
wio update pkg <directory>
```

## Build
Building a wio project is very easy and it comes with many configurations you can uses. Here are examples:

```bash
# build default target
wio build
```
```bash
# build specific target
wio build -target another1
```
```bash
# clean build files before building
wio build --clean
```
```bash
# if not in the project directory
wio build -dir project-dir
```

All of these commands above will build the project. Wio takes care of all the build files, dependency trees, etc. You just need to put files in proper project structure and all will be good. Wio uses `CMake` and `Make` to build the project. These build files are generated when `build` command is executed and hence the latest changes are picked up from `wio.yml` file. All the build files and internal helpers are stored in `.wio` folder. This folder should never be touched and modified since it can cause unexpected behaviour.

## Clean
Cleaning wio build files are easy and it also comes with few configurations you can uses. Here are examples:

```bash
# clean all the targets
wio clean
```
```bash
# clean a specific target
wio clean -target another1
```
```bash
# if not in the project directory
wio clean -dir project-dir
```

## Run
`run` command is a unique command in wio as it allows the user to do multiple things in one go. With `run` command, you can build, clean, and upload the project. In future releases a support for test command will be added as well. This is convinient for the user as the examples are below:

```bash
# normal build (default target)
wio run
```
```bash
# build with clean (default target)
wio run --clean
```
```bash
# build and upload (default target)
wio run --port somePort
```
```bash
# build and upload (specific target)
wio run -target someTarget -port somePort
```

As you can see above, upload is automatically triggered if port is provided. At the moment automcatic port detection does not work and user needs to provide the port to upload the hex file to.

## Package Manager
A goal of Wio has always been to make embedded development easier and modular. No more clone the repository and linking the code. This is achieved by wio package manager called `pac`. Wio is using `npm` package manager as a backend and a lot of warpper code has been written around it make it compatible with embedded development. Let's talk about command and it's features:

### Publish
```bash
# publish pkg from current directory
wio pac publish 
```

```bash
# publish pkg from another directory
wio pac publish -dir someDir
```

`publish` command is used when your `pkg` project is ready and now you want to publish it for other people to use. You can fill in all the meta information and other development information. wio will make a tarball and push the package to npm. **Note: test directory is not included in this process**. Multiple checks are made to make sure all the dependencies are valid for this package. Validity is checked as follows:
* If `vendor: true` for the dependency, it is not included as a dependency for `npm` package.
* A check is made to see if the dependency exists on npm backend.
    * Both name and version are verified.
* A check is made to see if that package is a valid wio project.

If the name of the package clashes with the name already on the server, an error message is thrown. Everytime a `publish` command is executed, you need to update the version for the project to publish it again next time. Before you can use `publish` command, you need to do the following:

```bash
# will create an account or log you in (ignore if already logged in)
npm addusr
```

### Get
`get` command is used to pull in all the dependencies for the project. When this command is executed, wio goes through all the listed dependencies in `wio.yml` file and pulls all of them. It also pulls the whole dependency tree that comes with these dependencies. In the build process, these dependencies are linked to compile the project. Note this applies to "remote" dependencies and not "vendor" dependencies. Vendor dependencies are discussed below. Example of `get` command

```bash
# pulls all the dependencies for the current directory
wio pac get
```
```bash
# pulls all the dependencies and clean the old ones
wio pac get --clean
```
```bash
# pulls all the dependencies for the provided directory
wio pac get -dir someDir
```

This command also acts as an update command as well. You can modify the versions of dependencies in `wio.yml` file and then run this command. This will automatically upgrade or downgrade the dependency and update the dependency tree.

### Vendor
This is a feature inspired by `golang`. If a vendor directory is created in any wio project, it is treated as a special directory for dependencies. All the dependencies placed in this folder are not managed by the package manager but rather they get precedence over the remote dependencies. This is useful when you do not want to publish a packages but still want to use it in other projects. During the build process, these packages will be automatically linked and compiled. **Note: You still need to specify these packages in `wio.yml` file as dependencies**. If they are not specified as dependencies, they will not be included in the build process.

## Wcosa -> Wio
This project used be known as `wcosa` but, our team had decided to completely rewamp the project
and hence change the name to `wio`. `wio` is a development environment for avr chips. It is 
more powerful than ever and is a direct replacement of `wcosa`. In future support for Arduino SDK
will be added but at the moment only `cosa` sdk is supported. In order to migrate your `wcosa` project,
follow the following steps:

Steps for `app` projects:
```bash
# go to project directory
cd <wcosa-directory>
```
```bash
# need to deleted to be replaced by wio .gitignore
rm -rf .gitignore
```
```bash
# initialize the wcosa project as wio project (update will add missing files)
wio update app . <board>
```
```bash
# remove wcosa specific directoris
rm -rf lib and wcosa

# remove ide (clion) files if there
rm -rf CMakeLists.txt CMakeListsPrivate.txt
```
```bash
# open config.json file and based on that update wio.yml file
# open pkglist.json file and based on that add dependencies to wio.yml file
# then remove the config.json file
rm -rf config.json pkglist.json
```

Steps for `pkg` projects:
```bash
# go to project directory
cd <wcosa-directory>
```
```bash
# need to deleted to be replaced by wio .gitignore
rm -rf .gitignore
```
```bash
# initialize the wcosa project as wio project (update will add missing files)
wio update pkg . <board>
```
```bash
# transfer public headers froms src to include
```
```bash
# remove wcosa specific directoris
rm -rf lib and wcosa

# remove ide (clion) files if there
rm -rf CMakeLists.txt CMakeListsPrivate.txt
```
```bash
# open config.json file and based on that update wio.yml file (add meta information for pkg as well)
# open pkglist.json file and based on that add dependencies to wio.yml file
# then remove the config.json file
rm -rf config.json pkglist.json
```
