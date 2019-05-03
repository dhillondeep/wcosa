package testing

// case 1: no scope values provided, should override with global
var noScopeValues = `
type: pkg

project:
  name: noScopeValues
  compile_options:
    flags: [flag1, flag2]
    definitions: [def1, def2]
    cxx_standard: c++11
    c_standard: c11
  package_options:
    header_only: true
    type: STATIC

targets:
  - name: main
    arguments: [arg1, arg2]
`

// case 2: some scope values provided, should add others from global
var someScopeValues = `
type: pkg

project:
  name: someScopeValues
  compile_options:
    flags: [flag1, flag2]
    definitions: [def1, def2]
    cxx_standard: c++11
    c_standard: c11
  package_options:
    header_only: true
    type: STATIC

targets:
  - name: main
    arguments: [arg1, arg2]
    compile_options:
      cxx_standard: c++17
      c_standard: c01
    package_options:
      type: SHARED
  - name: main2
    package_options:
      header_only: false
`

// case 3: array values are provided, append them with global
var arrayValuesProvided = `
type: app

project:
  name: arrayValuesProvided
  compile_options:
    flags: [flag1, flag2]
    definitions: [def1, def2]
    cxx_standard: c++11
    c_standard: c11

targets:
  - name: main
    arguments: [arg1, arg2]
    compile_options:
      flags: [flag3]
      definitions: [def3]
`

// case 4: app config warnings
var appConfigWarnings = `
type: app

project:
  name: appConfigWarnings
  compile_options:
    flags: [flag1, flag2]
    definitions: [def1, def2]
    cxx_standard: c++11
    c_standard: c11
  package_options:
    header_only: false

targets:
  - name: main
    executable_options:
      main_file: src/hello.cpp
    arguments: [arg1, arg2]
    compile_options:
      flags: [flag3]
      definitions: [def3]
    package_options:
      type: SHARED

tests:
  - name: main
    executable_options:
      main_file: src/hello
`

// case 5: pkg config warnings
var pkgConfigWarnings = `
type: pkg

project:
  name: pkgConfigWarnings
  compile_options:
    flags: [flag1, flag2]
    definitions: [def1, def2]
    cxx_standard: c++11
    c_standard: c11
  package_options:
    header_only: false

targets:
  - name: main
    executable_options:
      main_file: src/hello.cpp
    arguments: [arg1, arg2]
    compile_options:
      flags: [flag3]
      definitions: [def3]
    package_options:
      type: SHARED

tests:
  - name: main
    executable_options:
      main_file: src/hello.cpp
`

// case 6: convert a string to string array for certain fields
var stringToSliceFields = `
type: app

project:
  name: stringToSliceFields
  contributors: Jordan
  repository: repo
  compile_options:
    flags: flag1
    definitions: def1

variables: var1=10
arguments: Debug

targets:
  - name: main
    executable_options:
      toolchain: github.com/wio-pm/someToolchain@develop
      source: src
    arguments: NOP
    compile_options:
      flags: flag2
      definitions: def2
    linker_options:
      flags: link1

  - bru

tests:
  - name: main
    executable_options:
      toolchain: github.com/wio-pm/someToolchain
      source: test
    arguments: NOP
    compile_options:
      flags: flag2
      definitions: def2
    linker_options:
      flags: link1

dependencies:
  - github.com/random/dep

test_dependencies:
  - github.com/random/test_dep@develop
`

// case 7: convert a string to string array if separated by ,
var stringToSliceFieldsCommas = `
type: app

project:
  name: stringToSliceFields
  contributors: Jordan, Simon
  repository: repo, repo2
  compile_options:
    flags: flag1, flag2
    definitions: def1, def2

variables: var1=10 , var2=20
arguments: Debug, Holy=5

targets:
  - name: main
    executable_options:
      toolchain: someToolchain@default
      source: src, common, utils
    arguments: NOP, MOP
    compile_options:
      flags: flag2, flag4
      definitions: def2, def4
    linker_options:
      flags: link1, link2

tests:
  - name: main
    executable_options:
      toolchain: someToolchain@default
      source: test, utils
    arguments: NOP, MOP
    compile_options:
      flags: flag2, flag4
      definitions: def2, def4
    linker_options:
      flags: link1, link2
`

// case 8: random file content
var randomFileContent = `
type = pkg

hello
gg:
no way:

compile_options:
`

// case 9: valid file wrong schema
var invalidSchema = `
type: pkg

project:
  name: SampleProject
  repository:
    repo1: name1
  author: Deep
`

// case 10: unsupported tag
var unsupportedTag = `
type: pkg

project:
  type: Project
  name: SampleProject
  author: Deep
`

// case 11: hil usage
var hilUsage = `
type: app

variables:
  - PROJECT_NAME=HilUsage
  - AUTHOR=Test
  - SCRIPT1=somePathScript1
  - SCRIPT2=somePathScript2
  - TEST_MAIN_VISIBILITY=PUBLIC
  - DEP_ONE_REF=default
  - DEP_ONE_VISIBILITY=PRIVATE

arguments:
  - HOMEPAGE=homepage

scripts:
  begin: ${var.SCRIPT1}
  end: ${var.SCRIPT2}

project:
  name: ${var.PROJECT_NAME}
  author: ${var.AUTHOR}
  version: 0.0.1
  homepage: ${arg.HOMEPAGE}
  description: ${append(var.PROJECT_NAME, " project description")}

targets:
  - name: main
    executable_options:
      source: ${lower("SRC")}
      platform: '${var.PROJECT_NAME == "HilUsage" ? "native" : "windows"}'
      toolchain: '${arg.DEBUG == false ? "prod" : "debug"}@${arg.DEBUG == false ? "default" : "main"}'

    arguments:
      - DEBUG=false

tests:
  - name: main
    target_name: main
    target_arguments:
      - DEBUG=${lower("TRUE")}

    arguments:
      - VISIBILITY_CHECK=true

    linker_options:
      visibility: '${arg.VISIBILITY_CHECK == true ? var.TEST_MAIN_VISIBILITY : "SHARED"}'

dependencies:
  - name: gitlab.com/user/dependency1
    ref: ${var.DEP_ONE_REF}
    arguments:
      - DEBUG=true
    linker_options:
      visibility: ${var.DEP_ONE_VISIBILITY}

test_dependencies:
  - name: gitlab.com/user/dependency1
    ref: ${var.DEP_ONE_REF}
    arguments:
      - DEBUG=true
    linker_options:
      visibility: ${var.DEP_ONE_VISIBILITY}
`

// case 12: invalid hil usage and invalid version
var invalidHilUsage = `
type: pkg

variables:
  - VERSION=randomV0.1.0

scripts:
  begin: '${randomFunc(var.SCRIPT1)}'

project:
  name: ${${arg1.Blah}}
  version: ${var.VERSION}
  package_options:
    header_only: ${var.VERSION}
`

// case 13: script exec for fields
var configScriptExec = `
type: app

variables:
  - PROJECT_NAME=ExecUsage

project:
  name: $exec{out = "${var.PROJECT_NAME}"}
  description: '$exec{
    out = "${var.PROJECT_NAME}" + " project description"
  }'

targets:
  - name: main
    executable_options:
      source: '$exec {
        text := import("text")
        
        out = text.to_lower("SRC")
      }'

      platform: '$exec {
        if "${var.PROJECT_NAME}" == "ExecUsage" {
          out = "native"
        } else {
          out = "windows"
        }
      }'
`

// case 14: script exec invalid
var configScriptExecInvalid = `
type: pkg

variables:
  - PROJECT_NAME=ExecUsage

project:
  name: $exec{out = randomFunc("${var.PROJECT_NAME}")}
  version: $exec{out = RANDOM_VERSION}
  description: '$exec{
    out = "${${var.PROJECT_NAME}}" + " project description"
  }'
  package_options:
    header_only: $exec{out = ${var.PROJECT_NAME}}

tests:
  - name: main
    executable_options:
      source: '$exec {
        text := import("text")
        
        out = text.to_lower(src)
      }'
`

var createConfigAppWithToolchain = `type: app

project:
  name: %s
  version: 0.0.1

targets:
  - name: main
    executable_options:
      source: src
      main_file: %s
      platform: %s
      toolchain: %s

tests:
  - name: main
    executable_options:
      source: test
      platform: %s
      toolchain: %s

    target_name: main
`

var createConfigAppWithoutToolchain = `type: app

project:
  name: %s
  version: 0.0.1

targets:
  - name: main
    executable_options:
      source: src
      main_file: %s
      platform: %s

tests:
  - name: main
    executable_options:
      source: test
      platform: %s

    target_name: main
`

var createConfigPkgHeaderOnlyToolchain = `type: pkg

project:
  name: %s
  version: 0.0.1

  package_options:
    header_only: true

targets:
  - %s

tests:
  - name: %s
    executable_options:
      source: test
      platform: %s
      toolchain: %s

    target_name: %s
`

var createConfigPkgWithoutToolchain = `type: pkg

project:
  name: %s
  version: 0.0.1

targets:
  - %s

tests:
  - name: %s
    executable_options:
      source: test
      platform: %s

    target_name: %s
`

var createConfigPkgHeaderOnlyShared = `type: pkg

project:
  name: %s
  version: 0.0.1

  package_options:
    type: SHARED
    header_only: true

targets:
  - %s

tests:
  - name: %s
    executable_options:
      source: test
      platform: %s

    target_name: %s
`

var createConfigPkgShared = `type: pkg

project:
  name: %s
  version: 0.0.1

  package_options:
    type: SHARED

targets:
  - %s

tests:
  - name: %s
    executable_options:
      source: test
      platform: %s

    target_name: %s
`
