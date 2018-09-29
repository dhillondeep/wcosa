[![Build Status](https://travis-ci.org/wio/wio.svg?branch=dev)](https://travis-ci.org/wio/wio) [![npm](https://img.shields.io/npm/dw/wio.svg)](https://www.npmjs.com/package/wio) [![npm](https://img.shields.io/npm/v/wio.svg)](https://www.npmjs.com/package/wio) [![GitHub release](https://img.shields.io/github/release/wio/wio.svg)](https://github.com/wio/wio/releases) [![license](https://img.shields.io/github/license/wio/wio.svg)](https://github.com/wio/wio/blob/master/LICENSE)

**Quick Links:** [Docs](https://wio.github.io/docs) | [Packages](https://www.npmjs.com/search?q=wio%2C%20pkg) | [Waterloop](https://waterloop.ca)

![wio](https://wio.github.io/docs/_static/logo_black.png)

Wio is a development environment to create, build, test and upload C/C++ project. Wio supports `AVR` and `Native` platform but as the project grows, more platforms and frameworks will be added.

## Introduction
* [What is Wio?](https://wio.github.io/docs/wio/what-is-wio.html)

## Features
### Easier Development
* User worries about programming and wio takes care of build files
* No need to learn about vendor toolchains
* Support for multiple platforms and frameworks
* Highly customizable with the help of wio.yml config
* Build projects with `wio build`
* Execute projects with `wio run`
### Package Manager
* Code sharing has never been easier with wio packages
* User includes dependencies and wio takes care of dependency tree
* Use `wio publish` to publish a package
* Use `wio install <package>` to install a package
### Devices
* Uploading code to devices can be done with `wio run --port <port>`
* List devices connected to machine by using `wio devices list`
* Open a Serial monitor using `wio devices monitor`

## Installation
* [Linux](https://wio.github.io/docs/wio/install/linux.html)
* [MacOS](https://wio.github.io/docs/wio/install/macos.html)
* [Windows](https://wio.github.io/docs/wio/install/windows.html)

### Install from source
* Clone dev branch for latest features 
```bash
git clone --recurse-submodules https://github.com/wio/wio.git
```
* Install tools and build
```bash
# windows
wmake setup
wmake build

# macOs
wmake mac-setup
wmake build

# Linux
wmake linux-setup
wmake build
```



## Getting Started
* [Application Project](https://wio.github.io/docs/wio/getting-started/creating-app.html)
* [Package Project](https://wio.github.io/docs/wio/getting-started/creating-pkg.html)

## Contributing
If you are interested in working on wio, you can read [contribution document](https://github.com/wio/wio/blob/master/CONTRIBUTING.md) and add features/fixes.
