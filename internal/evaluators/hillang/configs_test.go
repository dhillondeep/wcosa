package hillang

import (
	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thoas/go-funk"
	"testing"
	"wio/internal/config"
)

type ConfigTestSuite struct {
	suite.Suite
	Variables   config.Variables
	Arguments   config.Arguments
	DefaultEval *hil.EvalConfig
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.Variables = config.Variables{
		config.VariableImpl{
			Name:  "variableOne",
			Value: "One",
		},
		config.VariableImpl{
			Name:  "variableTwo",
			Value: "Two",
		},
	}

	suite.Arguments = config.Arguments{
		config.ArgumentImpl{
			Name:  "argumentOne",
			Value: "One",
		},
		config.ArgumentImpl{
			Name:  "argumentHil",
			Value: "${2 + 2}",
		},
		config.ArgumentImpl{
			Name:  "argumentHilInvalid",
			Value: "${Random}",
		},
	}

	suite.DefaultEval = &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap:  map[string]ast.Variable{},
			FuncMap: getFunctions(),
		},
	}

	evalConfig = suite.DefaultEval
}

func assertVariables(t *testing.T, variables config.Variables, evalConfig *hil.EvalConfig) {
	for _, variable := range variables {
		key := "var." + variable.GetName()
		assert.True(t, funk.Contains(funk.Keys(evalConfig.GlobalScope.VarMap), key))
		assert.Equal(t, evalConfig.GlobalScope.VarMap[key].Value, variable.GetValue())
	}
}

func assertArguments(t *testing.T, arguments config.Arguments, evalConfig *hil.EvalConfig, isError bool) {
	for _, argument := range arguments {
		key := "arg." + argument.GetName()
		assert.True(t, funk.Contains(funk.Keys(evalConfig.GlobalScope.VarMap), key))

		value, err := argument.GetValue(evalConfig)

		if !isError {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}

		assert.Equal(t, evalConfig.GlobalScope.VarMap[key].Value, value)
	}
}

func (suite *ConfigTestSuite) TestInitialize() {
	// happy path - basic
	varsToUse := append(config.Variables{}, suite.Variables[0])
	argsToUse := append(config.Arguments{}, suite.Arguments[0])

	err := Initialize(varsToUse, argsToUse)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), evalConfig, nil)

	assertVariables(suite.T(), varsToUse, evalConfig)
	assertArguments(suite.T(), argsToUse, evalConfig, false)

	// happy path - hil eval arguments
	varsToUse = append(config.Variables{}, suite.Variables[0])
	argsToUse = append(config.Arguments{}, suite.Arguments[1])

	err = Initialize(varsToUse, argsToUse)
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), evalConfig, nil)

	assertVariables(suite.T(), varsToUse, evalConfig)
	assertArguments(suite.T(), argsToUse, evalConfig, false)

	// error - invalid hil
	varsToUse = append(config.Variables{}, suite.Variables[0])
	argsToUse = append(config.Arguments{}, suite.Arguments[2])

	err = Initialize(varsToUse, argsToUse)
	assert.Error(suite.T(), err)
}

func (suite *ConfigTestSuite) TestGetDefaultEvalConfig() {
	// happy path - default eval config
	returned := GetDefaultEvalConfig()
	assert.Equal(suite.T(), evalConfig, returned)

	// happy path - basic
	evalConfig = &hil.EvalConfig{
		GlobalScope:    nil,
		SemanticChecks: nil,
	}

	returned = GetDefaultEvalConfig()
	assert.Equal(suite.T(), evalConfig, returned)
}

func (suite *ConfigTestSuite) TestGetArgsEvalConfig() {
	varsToUse := append(config.Variables{}, suite.Variables[0])
	argsToUse := append(config.Arguments{}, suite.Arguments[0])

	// happy path - arguments are appended
	moreArgs := append(config.Arguments{}, suite.Arguments[1])

	_ = Initialize(varsToUse, argsToUse)

	newConfig, err := GetArgsEvalConfig(moreArgs, evalConfig)
	assert.NoError(suite.T(), err)

	assertVariables(suite.T(), varsToUse, newConfig)
	assertArguments(suite.T(), append(argsToUse, moreArgs...), newConfig, false)

	// happy path - argument is overridden
	moreArgs = append(config.Arguments{}, config.ArgumentImpl{
		Name:  suite.Arguments[0].GetName(),
		Value: "NewTwo",
	})

	newConfig, err = GetArgsEvalConfig(moreArgs, evalConfig)
	assert.NoError(suite.T(), err)

	assertVariables(suite.T(), varsToUse, newConfig)
	assertArguments(suite.T(), moreArgs, newConfig, false)

	overriddenVal, err := moreArgs[0].GetValue(evalConfig)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), overriddenVal, newConfig.GlobalScope.VarMap["arg."+moreArgs[0].GetName()].Value)

	// happy path - two evalConfigs are different
	assert.NotEqual(suite.T(), newConfig, evalConfig)
	assert.NotEqual(suite.T(), newConfig.GlobalScope.VarMap, evalConfig.GlobalScope.VarMap)

	// error - hil eval fails
	moreArgs = append(config.Arguments{}, suite.Arguments[2])
	newConfig, err = GetArgsEvalConfig(moreArgs, evalConfig)
	assert.Error(suite.T(), err)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
