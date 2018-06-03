// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for each of them

package cmake

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/parsers"
    "wio/cmd/wio/utils"
    "io/ioutil"
    "path/filepath"
    "strings"
    "errors"
)

var configFileName = "wio.yml"
var vendorPkgFolder = "vendor"
var remotePkgFolder = "node_modules"

// Parses Package based on the path provides and its dependencies and it modifies dependency tree to store the data
func parsePackage(rootPackagesPath string, currPackagePath string, depTree *parsers.DependencyTree,
    parentDependencyData *types.DependencyTag, hash string) (error) {

    configFile := currPackagePath + io.Sep + configFileName

    var config types.PkgConfig
    var pkgName string

    dependencyTag := parsers.DependencyTag{}
    dependencies := types.DependenciesTag{}

    // check if wio.yml file exists for the dependency and based on that gather the name and its dependencies
    if utils.PathExists(configFile) {
        if err := io.NormalIO.ParseYml(configFile, &config); err != nil {
            return err
        }

        pkgName = config.MainTag.Name
        dependencyTag.Compile_flags = config.MainTag.CompileFlags
        dependencies = config.DependenciesTag
    } else {
        pkgName = filepath.Base(currPackagePath)
        return errors.New(pkgName + " dependency is not a wio package")
    }

    if len(hash) <= 0 {
        hash += pkgName
    } else {
        hash += "__" + pkgName
    }

    dependencyTag.Name = pkgName
    dependencyTag.Hash = hash
    dependencyTag.Path = currPackagePath

    depTree.Config = dependencyTag
    depTree.Child = make(map[string]*parsers.DependencyTree, 0)

    // parse vendor dependencies
    vendorPath := currPackagePath + io.Sep + "vendor"
    if utils.PathExists(vendorPath) {
        if dirs, err := ioutil.ReadDir(vendorPath); err != nil {
            return err
        } else if len(dirs) > 0 {
            // parse vendor dependencies
            for depName, depValue := range dependencies {
                // only parse vendor packages
                if !depValue.Vendor {
                    continue
                }

                depPath := vendorPath + io.Sep + depName
                currTree := parsers.DependencyTree{}

                if err := parsePackage(vendorPath, depPath, &currTree, depValue, hash);
                    err != nil {
                    return err
                }

                if _, ok := depTree.Child[currTree.Config.Name]; ok {
                    continue
                } else {
                    depTree.Child[currTree.Config.Name] = &currTree
                }
            }
        }
    }

    // parse remote dependencies
    for depName, depValue := range dependencies {
        // only parse non vendor packages
        if depValue.Vendor {
            continue
        }

        depPath := rootPackagesPath + io.Sep + depName
        currTree := parsers.DependencyTree{}

        if err := parsePackage(rootPackagesPath, depPath, &currTree, depValue, hash);
            err != nil {
            return err
        }

        if _, ok := depTree.Child[currTree.Config.Name]; ok {
            continue
        } else {
            depTree.Child[currTree.Config.Name] = &currTree
        }
    }

    if parentDependencyData != nil {
        depTree.Config.Compile_flags = utils.AppendIfMissing(parentDependencyData.CompileFlags,
            depTree.Config.Compile_flags)
    }

    return nil
}

// Parses all the packages recursively with the help of ParsePackage function
func parsePackages(packagesPath string, depTrees map[string]*parsers.DependencyTree, dependencies types.DependenciesTag,
    hash string, vendor bool) (map[string]*parsers.DependencyTree, error) {

    if !utils.PathExists(packagesPath) {
        return depTrees, nil
    }

    if dirs, err := ioutil.ReadDir(packagesPath); err != nil {
        return depTrees, err
    } else if len(dirs) > 0 {
        for depName, depValue := range dependencies {
            // only parse vendor packages
            if vendor && depValue.Vendor {
                continue
            }

            depPath := packagesPath + io.Sep + depName

            // check if dependency exists
            if !utils.PathExists(depPath) {
                return nil, errors.New(depName + " dependency does not exist. use wio pac get to pull the dependency")
            }

            currTree := parsers.DependencyTree{}

            if err := parsePackage(packagesPath, depPath, &currTree, depValue, hash);
                err != nil {
                return depTrees, err
            }

            if _, ok := depTrees[currTree.Config.Name]; ok {
                continue
            } else {
                depTrees[currTree.Config.Name] = &currTree
            }
        }
    }

    return depTrees, nil
}

// Parses all the packages and their dependencies and creates pkg.lock file
func createDependencyTree(projectPath string, dependencies types.DependenciesTag) (map[string]*parsers.DependencyTree, error) {
    wioPath := projectPath + io.Sep + ".wio"
    librariesLocalPath := projectPath + io.Sep + vendorPkgFolder
    librariesRemotePath := wioPath + io.Sep + remotePkgFolder

    dependencyTrees := make(map[string]*parsers.DependencyTree, 0)

    // parse all the libraries and their dependencies in vendor folder
    dependencyTrees, err := parsePackages(librariesLocalPath, dependencyTrees, dependencies, "", true)
    if err != nil {
        return nil, err
    }

    // parse all the libraries and their dependencies in node_modules folder
    dependencyTrees, err = parsePackages(librariesRemotePath, dependencyTrees, dependencies, "", false)
    if err != nil {
        return nil, err
    }

    // return Dependency Tree
    return dependencyTrees, nil
}

func buildCMakeString(templateString string, finalString string, depTree *parsers.DependencyTree) (string) {
    currString := templateString

    currString = strings.Replace(currString, "{{DEPENDENCY_PATH}}", filepath.ToSlash(depTree.Config.Path), -1)
    currString = strings.Replace(currString, "{{DEPENDENCY_NAME}}", depTree.Config.Hash, -1)
    currString = strings.Replace(currString, "{{DEPENDENCY_COMPILE_FLAGS}}",
        strings.Join(depTree.Config.Compile_flags, " "), -1)
    currString += "\n"

    for dep := range depTree.Child {
        currLinkString := "target_link_libraries({{DEPENDENCY_NAME}} PUBLIC {{DEP_NAME}})"
        currLinkString = strings.Replace(currLinkString, "{{DEPENDENCY_NAME}}", depTree.Config.Hash, -1)
        currLinkString = strings.Replace(currLinkString, "{{DEP_NAME}}", depTree.Child[dep].Config.Hash, -1)

        currString += currLinkString + "\n"
    }

    finalString += currString + "\n"

    for dep := range depTree.Child {
        finalString = buildCMakeString(templateString, finalString, depTree.Child[dep])
    }

    return finalString
}

// Given all the dependency tree data, it will create cmake files for each dependency
func createDependencyCMakeString(depTree *parsers.DependencyTree) (string, error) {

    dependenciesTemplate, err := io.AssetIO.ReadFile("templates" + io.Sep + "cmake" + io.Sep +
        "dependencies.txt")
    if err != nil {
        return "", err
    }

    dependencyBuildString := buildCMakeString(string(dependenciesTemplate), "", depTree)

    return dependencyBuildString, nil
}

// This parses dependencies and creates cmake files for all of these. It does this for a target that is provided
// to it.
func ParseDepsAndCreateCMake(projectPath string, dependencies types.DependenciesTag) (map[string]*parsers.DependencyTree, error) {

    // create dependency tree
    dependencyTree, err := createDependencyTree(projectPath, dependencies)
    if err != nil {
        return nil, err
    }

    finalBuildString := ""

    for tree := range dependencyTree {
        if str, err := createDependencyCMakeString(dependencyTree[tree]);
            err != nil {
            return nil, err
        } else {
            finalBuildString += str
        }
    }

    if err := io.NormalIO.WriteFile(projectPath+io.Sep+".wio"+io.Sep+"build"+io.Sep+"dependencies.cmake",
        []byte(finalBuildString)); err != nil {
        return nil, err
    }

    return dependencyTree, nil
}
