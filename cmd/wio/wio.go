package main

import (
	"fmt"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"wio/internal/config"
	"wio/internal/hillang"
	"wio/pkg/sys"
)

func readFileText(filePath string) string {
	data, err := sys.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func main() {
	// config.CreateConfig(templates.ProjectCreation{
	// 	Type: constants.APP,
	// 	ProjectName: "sampleApp",
	// 	ProjectPath: "/Users/deep/Development/waterloop/projects/wio/tmp",
	// 	Platform: "native",
	// 	// Toolchain: "s",
	// 	MainFile: "src/main.cpp",
	// })
	//
	// config.CreateConfig(templates.ProjectCreation{
	// 	Type: constants.PKG,
	// 	ProjectName: "samplePkg",
	// 	ProjectPath: "/Users/deep/Development/waterloop/projects/wio/tmp",
	// 	Platform: "native",
	// 	Toolchain: "gg",
	// 	PkgType: "static",
	// 	HeaderOnly: true,
	// })

	config2, err := config.ReadConfig("/Users/deep/Development/waterloop/projects/wio/tmp")
	if err != nil {
		panic(err)
	}

	hillang.Initialize(config2.GetVariables(), config2.GetArguments())
	evalConfig := hillang.GetDefaultEvalConfig()

	author, err := config2.GetProject().GetAuthor(evalConfig)
	if err != nil {
		panic(err)
	}

	version, err := config2.GetProject().GetVersion(evalConfig)
	if err != nil {
		panic(err)
	}

	description, err := config2.GetProject().GetDescription(evalConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println(author)
	fmt.Println(version)
	fmt.Println(description)

	newConfig, err := hillang.GetArgsEvalConfig(config2.GetTests()["main"].GetArguments())
	if err != nil {
		panic(err)
	}

	targetName, err := config2.GetTests()["main"].GetTargetName(newConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println(targetName)

	env := vm.NewEnv()

	err = env.Define("println", fmt.Println)
	if err != nil {
		panic(err)
	}

	err = env.Define("readFile", readFileText)
	if err != nil {
		panic(err)
	}

	smt, err := parser.ParseSrc(description)
	if err != nil {
		panic(err)
	}

	val, err := env.Run(smt)
	fmt.Println("Description", val)
	if err != nil {
		panic(err)
	}
}
