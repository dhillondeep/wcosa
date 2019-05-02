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
			Value: "${2 + 2}",
		},
		config.VariableImpl{
			Name:  "variableTwo",
			Value: "${Random}",
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

func assertVariables(t *testing.T, variables config.Variables, evalConfig *hil.EvalConfig, isError bool) {
	for _, variable := range variables {
		key := "var." + variable.GetName()
		assert.True(t, funk.Contains(funk.Keys(evalConfig.GlobalScope.VarMap), key))

		value, err := variable.GetValue(evalConfig)

		if !isError {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}

		assert.Equal(t, evalConfig.GlobalScope.VarMap[key].Value, value)
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
	suite.T().Run("happy path - raw text values", func(t *testing.T) {
		varsToUse := append(config.Variables{}, suite.Variables[0])
		argsToUse := append(config.Arguments{}, suite.Arguments[0])

		err := Initialize(varsToUse, argsToUse)
		assert.NoError(suite.T(), err)
		assert.NotEqual(suite.T(), evalConfig, nil)

		assertVariables(suite.T(), varsToUse, evalConfig, false)
		assertArguments(suite.T(), argsToUse, evalConfig, false)
	})

	suite.T().Run("happy path - hil eval variables and arguments", func(t *testing.T) {
		varsToUse := append(config.Variables{}, suite.Variables[1])
		argsToUse := append(config.Arguments{}, suite.Arguments[1])

		err := Initialize(varsToUse, argsToUse)
		assert.NoError(suite.T(), err)
		assert.NotEqual(suite.T(), evalConfig, nil)

		assertVariables(suite.T(), varsToUse, evalConfig, false)
		assertArguments(suite.T(), argsToUse, evalConfig, false)
	})

	suite.T().Run("wrong path - invalid hil for variables", func(t *testing.T) {
		varsToUse := append(config.Variables{}, suite.Variables[2])
		argsToUse := append(config.Arguments{}, suite.Arguments[0])

		err := Initialize(varsToUse, argsToUse)
		assert.Error(suite.T(), err)

	})

	suite.T().Run("wrong path - invalid hil for arguments", func(t *testing.T) {
		varsToUse := append(config.Variables{}, suite.Variables[0])
		argsToUse := append(config.Arguments{}, suite.Arguments[2])

		err := Initialize(varsToUse, argsToUse)
		assert.Error(suite.T(), err)
	})
}

func (suite *ConfigTestSuite) TestGetDefaultEvalConfig() {
	suite.T().Run("happy path - default config", func(t *testing.T) {
		returned := GetDefaultEvalConfig()
		assert.Equal(suite.T(), evalConfig, returned)
	})

	suite.T().Run("happy path - basic config", func(t *testing.T) {
		evalConfig = &hil.EvalConfig{
			GlobalScope:    nil,
			SemanticChecks: nil,
		}

		returned := GetDefaultEvalConfig()
		assert.Equal(suite.T(), evalConfig, returned)
	})
}

func (suite *ConfigTestSuite) TestGetArgsEvalConfig() {
	varsToUse := append(config.Variables{}, suite.Variables[0])
	argsToUse := append(config.Arguments{}, suite.Arguments[0])

	suite.T().Run("happy path - new arguments are appended", func(t *testing.T) {
		moreArgs := append(config.Arguments{}, suite.Arguments[1])

		_ = Initialize(varsToUse, argsToUse)

		newConfig, err := GetArgsEvalConfig(moreArgs, evalConfig)
		assert.NoError(t, err)

		assertVariables(t, varsToUse, newConfig, false)
		assertArguments(t, append(argsToUse, moreArgs...), newConfig, false)

	})

	suite.T().Run("happy path - argument is overridden", func(t *testing.T) {
		moreArgs := append(config.Arguments{}, config.ArgumentImpl{
			Name:  suite.Arguments[0].GetName(),
			Value: "NewTwo",
		})

		newConfig, err := GetArgsEvalConfig(moreArgs, evalConfig)
		assert.NoError(suite.T(), err)

		assertVariables(t, varsToUse, newConfig, false)
		assertArguments(t, moreArgs, newConfig, false)

		overriddenVal, err := moreArgs[0].GetValue(evalConfig)
		assert.NoError(t, err)

		assert.Equal(t, overriddenVal, newConfig.GlobalScope.VarMap["arg."+moreArgs[0].GetName()].Value)

	})

	suite.T().Run("happy path - two evalConfigs are different", func(t *testing.T) {
		moreArgs := append(config.Arguments{}, config.ArgumentImpl{
			Name:  suite.Arguments[0].GetName(),
			Value: "NewTwo",
		})

		newConfig, err := GetArgsEvalConfig(moreArgs, evalConfig)
		assert.NoError(t, err)

		// happy path - two evalConfigs are different
		assert.NotEqual(t, newConfig, evalConfig)
		assert.NotEqual(t, newConfig.GlobalScope.VarMap, evalConfig.GlobalScope.VarMap)

	})

	suite.T().Run("wrong path - hil eval fails for argument", func(t *testing.T) {
		moreArgs := append(config.Arguments{}, suite.Arguments[2])
		_, err := GetArgsEvalConfig(moreArgs, evalConfig)
		assert.Error(t, err)
	})
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
