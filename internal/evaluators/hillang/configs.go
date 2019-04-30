package hillang

import (
	"fmt"
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"wio/internal/config"
	"wio/internal/constants"
)

var globalArguments config.Arguments
var globalVariables config.Variables
var variablesMap = map[string]*config.Variable{}
var argsMap = map[string]*config.Argument{}

var evalConfig = &hil.EvalConfig{
	GlobalScope: &ast.BasicScope{
		VarMap:  map[string]ast.Variable{},
		FuncMap: getFunctions(),
	},
}

// copyEvalConfig is a helper function to deep copy one config to another
func copyEvalConfig(evalConfig *hil.EvalConfig) *hil.EvalConfig {
	newVars := map[string]ast.Variable{}

	for index, element := range evalConfig.GlobalScope.VarMap {
		newVars[index] = element
	}

	return &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			FuncMap: evalConfig.GlobalScope.FuncMap,
			VarMap:  newVars,
		},
		SemanticChecks: evalConfig.SemanticChecks,
	}
}

// GetDefaultEvalConfig provides default config. The config must be initialized in order to have
// proper variables and arguments
func GetDefaultEvalConfig() *hil.EvalConfig {
	return evalConfig
}

// GetArgsEvalConfig provides evalConfig with scope specific arguments
func GetArgsEvalConfig(arguments config.Arguments, evalConfig *hil.EvalConfig) (*hil.EvalConfig, error) {
	newConfig := copyEvalConfig(evalConfig)
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

// Initialize initializes variables and arguments to be used for hil parsing
func Initialize(variables config.Variables, arguments config.Arguments) error {
	globalVariables = variables
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
