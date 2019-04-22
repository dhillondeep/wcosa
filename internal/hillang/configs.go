package hillang

import (
	"fmt"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"wio/internal/config"
	"wio/internal/constants"
)

var globalArguments config.Arguments
var variablesMap = map[string]*config.Variable{}
var argsMap = map[string]*config.Argument{}

var evalConfig = &hil.EvalConfig{
	GlobalScope: &ast.BasicScope{
		VarMap:  map[string]ast.Variable{},
		FuncMap: getFunctions(),
	},
}

func GetDefaultEvalConfig() *hil.EvalConfig {
	return evalConfig
}

func GetArgsEvalConfig(arguments config.Arguments) (*hil.EvalConfig, error) {
	newConfig := evalConfig
	argsMap = map[string]*config.Argument{}

	for _, globalArg := range globalArguments {
		argsMap[globalArg.GetName()] = &globalArg
	}

	for _, givenArg := range arguments {
		argsMap[givenArg.GetName()] = &givenArg
	}

	for argName, arg := range argsMap {
		argumentValue, err := (*arg).GetValue(evalConfig)
		if err != nil {
			return nil, err
		}

		newConfig.GlobalScope.VarMap[fmt.Sprintf("%s.%s", constants.ARG, argName)] = ast.Variable{
			Type:  ast.TypeString,
			Value: argumentValue,
		}
	}

	return newConfig, nil
}

func Initialize(variables config.Variables, arguments config.Arguments) error {
	for _, variable := range variables {
		varName := variable.GetName()
		variablesMap[varName] = &variable
		evalConfig.GlobalScope.VarMap[fmt.Sprintf("%s.%s", constants.VAR, varName)] = ast.Variable{
			Type:  ast.TypeString,
			Value: variable.GetValue(),
		}
	}

	globalArguments = arguments
	for _, argument := range arguments {
		argumentValue, err := argument.GetValue(evalConfig)
		if err != nil {
			return err
		}

		argName := argument.GetName()
		argsMap[argName] = &argument
		evalConfig.GlobalScope.VarMap[fmt.Sprintf("%s.%s", constants.ARG, argName)] = ast.Variable{
			Type:  ast.TypeString,
			Value: argumentValue,
		}
	}

	return nil
}
