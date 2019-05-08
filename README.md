[![Build Status](https://travis-ci.org/wio/wio.svg?branch=dev)](https://travis-ci.org/wio/wio) [![Coverage Status](https://coveralls.io/repos/github/wio/wio/badge.svg?branch=develop-0.1.0)](https://coveralls.io/github/wio/wio?branch=develop-0.1.0) [![license](https://img.shields.io/github/license/wio/wio.svg)](https://github.com/wio/wio/blob/master/LICENSE)

**Quick Links:** [Waterloop](https://waterloop.ca)

![wio](https://wio.github.io/docs/_static/logo_black.png)

Wio is a development tool to create, build, and test C/C++ project. The idea behind this project is to
simplify development process for complex projects and for people who are not much familiar with C/C++ build tools. 
Wio uses config file (yaml) and cmake to provide a simple platform.

Wio has been in development for about a year and  many versions have been released https://github.com/wio/wio/releases. 
Throughout the process of building and adding new features, codebase has become unmanageable. This primarily has to do with the
design and the development process. Using all the leanings, good design practices, and growth of go, wio is being
re-developed. `wio v0.1.0` onwards will have following features on top of current features:
* Full template support
* Variables and arguments
* Ability to execute scripts
* Full support for cmake project, shared and static libraries
* Native support for testing and testing dependencies
* Easier process to support more platforms and toolchains
* Better configuration style
* Full testing and coverage of the code

With `wio v0.1.0` onwards, an example application will be as simple as:
```yaml
type: app

project:
  name: exampleProject
  version: 0.0.1
  
targets:
  - name: main
    executable_options:
      source: src
      platform: native
```

and package as simple as:
```yaml
type: pkg

project:
  name: examplePackage
  version: 0.0.1
  
targets:
  - main
```

Development for `wio v0.1.*` is being done in branch `develop-0.1.0`. The plan is to make necessary features available
on the rolling basis. If you are interested in using `wio v0.9.0` or below, please checkout the `master` code branch
https://github.com/wio/wio/tree/master.
 
## Contributing
If you are interested in working on wio, you can read [contribution document](https://github.com/wio/wio/blob/develop-0.1.0/CONTRIBUTING.md) and add features/fixes.
