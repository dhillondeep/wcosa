package testing

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"wio/internal/config"
	"wio/internal/constants"
	"wio/internal/evaluators/hillang"
	"wio/pkg/sys"
	"wio/templates"
)

const (
	noScopeValuesDir             = "/noScopeValues"
	someScopeValuesDir           = "/someScopeValues"
	arrayValuesProvidedDir       = "/arrayValuesProvided"
	appConfigWarningsDir         = "/appConfigWarnings"
	pkgConfigWarningsDir         = "/pkgConfigWarnings"
	stringToSliceFieldsDir       = "/stringToSliceFields"
	stringToSliceFieldsCommasDir = "/stringToSliceFieldsCommas"
	randomFileContentDir         = "/randomFileContent"
	invalidSchemaDir             = "/invalidSchema"
	unsupportedTagDir            = "/unsupportedTag"
	hilUsageDir                  = "/hilUsage"
	invalidHilUsageDir           = "/invalidHilUsage"
	configScriptExecDir          = "/configScriptEval"
	configScriptExecInvalidDir   = "/configScriptExecInvalid"
)

type ConfigTestSuite struct {
	suite.Suite
	config *hil.EvalConfig
}

func (suite *ConfigTestSuite) SetupTest() {
	sys.SetFileSystem(afero.NewMemMapFs())

	suite.config = hillang.GetDefaultEvalConfig()

	// save configs
	err := sys.WriteFile(noScopeValuesDir+sys.GetSeparator()+constants.WioConfigFile, []byte(noScopeValues))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(someScopeValuesDir+sys.GetSeparator()+constants.WioConfigFile, []byte(someScopeValues))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(arrayValuesProvidedDir+sys.GetSeparator()+constants.WioConfigFile, []byte(arrayValuesProvided))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(appConfigWarningsDir+sys.GetSeparator()+constants.WioConfigFile, []byte(appConfigWarnings))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(pkgConfigWarningsDir+sys.GetSeparator()+constants.WioConfigFile, []byte(pkgConfigWarnings))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(stringToSliceFieldsDir+sys.GetSeparator()+constants.WioConfigFile, []byte(stringToSliceFields))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(stringToSliceFieldsCommasDir+sys.GetSeparator()+constants.WioConfigFile,
		[]byte(stringToSliceFieldsCommas))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(randomFileContentDir+sys.GetSeparator()+constants.WioConfigFile, []byte(randomFileContent))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(invalidSchemaDir+sys.GetSeparator()+constants.WioConfigFile, []byte(invalidSchema))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(unsupportedTagDir+sys.GetSeparator()+constants.WioConfigFile, []byte(unsupportedTag))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(hilUsageDir+sys.GetSeparator()+constants.WioConfigFile, []byte(hilUsage))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(invalidHilUsageDir+sys.GetSeparator()+constants.WioConfigFile, []byte(invalidHilUsage))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(configScriptExecDir+sys.GetSeparator()+constants.WioConfigFile, []byte(configScriptExec))
	require.NoError(suite.T(), err)

	err = sys.WriteFile(configScriptExecInvalidDir+sys.GetSeparator()+constants.WioConfigFile,
		[]byte(configScriptExecInvalid))
	require.NoError(suite.T(), err)
}

func (suite *ConfigTestSuite) TestReadConfig() {
	suite.T().Run("Happy path - no scope values provided, should override with global", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(noScopeValuesDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, project.GetProject().GetPackageOptions(), project.GetTargets()["main"].GetPackageOptions())
		require.Equal(t, project.GetProject().GetCompileOptions(), project.GetTargets()["main"].GetCompileOptions())
	})

	suite.T().Run("Happy path - some scope values provided, should add others from global", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(someScopeValuesDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		expectedb, err := project.GetProject().GetPackageOptions().IsHeaderOnly(suite.config)
		require.NoError(t, err)
		actualb, err := project.GetTargets()["main"].GetPackageOptions().IsHeaderOnly(suite.config)
		require.NoError(t, err)

		require.Equal(t, expectedb, actualb)

		actualb, err = project.GetTargets()["main2"].GetPackageOptions().IsHeaderOnly(suite.config)
		require.NoError(t, err)

		require.Equal(t, false, actualb)

		actual, err := project.GetTargets()["main"].GetPackageOptions().GetPackageType(suite.config)
		require.NoError(t, err)
		require.Equal(t, "SHARED", actual)

		require.Equal(t, project.GetProject().GetCompileOptions().GetFlags(),
			project.GetTargets()["main"].GetCompileOptions().GetFlags())
		require.Equal(t, project.GetProject().GetCompileOptions().GetDefinitions(),
			project.GetTargets()["main"].GetCompileOptions().GetDefinitions())

		actual, err = project.GetTargets()["main"].GetCompileOptions().GetCXXStandard(suite.config)
		require.NoError(t, err)
		require.Equal(t, "c++17", actual)

		actual, err = project.GetTargets()["main"].GetCompileOptions().GetCStandard(suite.config)
		require.NoError(t, err)
		require.Equal(t, "c01", actual)
	})

	suite.T().Run("Happy path - array values are provided, append them with global", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(arrayValuesProvidedDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		globalFlags := project.GetProject().GetCompileOptions().GetFlags()
		targetFlags := project.GetTargets()["main"].GetCompileOptions().GetFlags()

		require.Equal(t, append(config.Flags{config.ExpressionImpl{Value: "flag3"}}, globalFlags...), targetFlags)

		globalDefinitions := project.GetProject().GetCompileOptions().GetDefinitions()
		targetDefinitions := project.GetTargets()["main"].GetCompileOptions().GetDefinitions()

		require.Equal(t, append(config.Definitions{config.ExpressionImpl{Value: "def3"}}, globalDefinitions...),
			targetDefinitions)

		globalCXX, err := project.GetProject().GetCompileOptions().GetCXXStandard(suite.config)
		require.NoError(t, err)
		globalC, err := project.GetProject().GetCompileOptions().GetCStandard(suite.config)
		require.NoError(t, err)

		targetCXX, err := project.GetTargets()["main"].GetCompileOptions().GetCXXStandard(suite.config)
		require.NoError(t, err)
		targetC, err := project.GetTargets()["main"].GetCompileOptions().GetCStandard(suite.config)
		require.NoError(t, err)

		require.Equal(t, globalCXX, targetCXX)
		require.Equal(t, globalC, targetC)
	})

	suite.T().Run("Warnings - app config warnings", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(appConfigWarningsDir)
		require.NoError(t, err)

		require.Equal(t, "app", project.GetType())

		require.Equal(t, 3, len(warnings))

		require.Equal(t, nil, project.GetProject().GetPackageOptions())
		require.Equal(t, nil, project.GetTargets()["main"].GetPackageOptions())

		file, err := project.GetTests()["main"].GetExecutableOptions().GetMainFile(suite.config)
		require.NoError(t, err)

		require.Equal(t, "", file)
	})

	suite.T().Run("Warnings - pkg config warnings", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(pkgConfigWarningsDir)
		require.NoError(t, err)

		require.Equal(t, "pkg", project.GetType())

		require.Equal(t, 2, len(warnings))

		require.Equal(t, nil, project.GetTargets()["main"].GetExecutableOptions())

		file, err := project.GetTests()["main"].GetExecutableOptions().GetMainFile(suite.config)
		require.NoError(t, err)

		require.Equal(t, "", file)
	})

	suite.T().Run("Happy path - convert a string to string array for certain fields", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(stringToSliceFieldsDir)
		require.NoError(t, err)

		require.Equal(t, 0, len(warnings))

		require.IsType(t, config.Contributors{}, project.GetProject().GetContributors())
		require.Equal(t, config.Contributors{config.ExpressionImpl{Value: "Jordan"}},
			project.GetProject().GetContributors())

		require.IsType(t, config.Repositories{}, project.GetProject().GetRepository())
		require.Equal(t, config.Repositories{config.ExpressionImpl{Value: "repo"}},
			project.GetProject().GetRepository())

		require.IsType(t, config.Flags{}, project.GetProject().GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag1"}},
			project.GetProject().GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, project.GetProject().GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def1"}},
			project.GetProject().GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Variables{}, project.GetVariables())
		require.Equal(t, config.Variables{config.VariableImpl{Name: "var1", Value: "10"}}, project.GetVariables())

		require.IsType(t, config.Arguments{}, project.GetArguments())
		require.Equal(t, config.Arguments{config.ArgumentImpl{Name: "Debug", Value: ""}}, project.GetArguments())

		targetToTest := project.GetTargets()["main"]

		require.IsType(t, config.Sources{}, targetToTest.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "src"}},
			targetToTest.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, targetToTest.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "someToolchain", Ref: "default"},
			targetToTest.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, targetToTest.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}, config.ExpressionImpl{Value: "flag1"}},
			targetToTest.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, targetToTest.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}, config.ExpressionImpl{Value: "def1"}},
			targetToTest.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, targetToTest.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"}}, targetToTest.GetLinkerOptions().GetFlags())

		testToTest := project.GetTests()["main"]

		require.IsType(t, config.Sources{}, testToTest.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "test"}},
			testToTest.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, testToTest.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "someToolchain", Ref: "default"},
			testToTest.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, testToTest.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}},
			testToTest.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, testToTest.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}},
			testToTest.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, testToTest.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"}}, testToTest.GetLinkerOptions().GetFlags())

		require.IsType(t, config.Dependencies{}, project.GetDependencies())
		require.IsType(t, &config.DependencyImpl{}, project.GetDependencies()["somedep"])
		require.Equal(t, &config.DependencyImpl{Ref: "default"}, project.GetDependencies()["somedep"])

		require.IsType(t, config.Dependencies{}, project.GetTestDependencies())
		require.IsType(t, &config.DependencyImpl{}, project.GetTestDependencies()["somedep"])
		require.Equal(t, &config.DependencyImpl{Ref: "default"}, project.GetTestDependencies()["somedep"])
	})

	suite.T().Run("Happy path - convert a string to string array if separated by ,", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(stringToSliceFieldsCommasDir)
		require.NoError(t, err)

		require.Equal(t, 0, len(warnings))

		require.IsType(t, config.Contributors{}, project.GetProject().GetContributors())
		require.Equal(t, config.Contributors{config.ExpressionImpl{Value: "Jordan"},
			config.ExpressionImpl{Value: "Simon"}}, project.GetProject().GetContributors())

		require.IsType(t, config.Repositories{}, project.GetProject().GetRepository())
		require.Equal(t, config.Repositories{config.ExpressionImpl{Value: "repo"},
			config.ExpressionImpl{Value: "repo2"}}, project.GetProject().GetRepository())

		require.IsType(t, config.Flags{}, project.GetProject().GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag1"}, config.ExpressionImpl{Value: "flag2"}},
			project.GetProject().GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, project.GetProject().GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def1"}, config.ExpressionImpl{Value: "def2"}},
			project.GetProject().GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Variables{}, project.GetVariables())
		require.Equal(t, config.Variables{config.VariableImpl{Name: "var1", Value: "10"},
			config.VariableImpl{Name: "var2", Value: "20"}}, project.GetVariables())

		require.IsType(t, config.Arguments{}, project.GetArguments())
		require.Equal(t, config.Arguments{config.ArgumentImpl{Name: "Debug", Value: ""},
			config.ArgumentImpl{Name: "Holy", Value: "5"}}, project.GetArguments())

		targetToTest := project.GetTargets()["main"]

		require.IsType(t, config.Sources{}, targetToTest.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "src"},
			config.ExpressionImpl{Value: "common"}, config.ExpressionImpl{Value: "utils"}},
			targetToTest.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, targetToTest.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "someToolchain", Ref: "default"},
			targetToTest.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, targetToTest.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}, config.ExpressionImpl{Value: "flag4"},
			config.ExpressionImpl{Value: "flag1"}, config.ExpressionImpl{Value: "flag2"}},
			targetToTest.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, targetToTest.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}, config.ExpressionImpl{Value: "def4"},
			config.ExpressionImpl{Value: "def1"}, config.ExpressionImpl{Value: "def2"}},
			targetToTest.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, targetToTest.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"},
			config.ExpressionImpl{Value: "link2"}}, targetToTest.GetLinkerOptions().GetFlags())

		testToTest := project.GetTests()["main"]

		require.IsType(t, config.Sources{}, testToTest.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "test"}, config.ExpressionImpl{Value: "utils"}},
			testToTest.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, testToTest.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "someToolchain", Ref: "default"},
			testToTest.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, testToTest.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}, config.ExpressionImpl{Value: "flag4"}},
			testToTest.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, testToTest.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}, config.ExpressionImpl{Value: "def4"}},
			testToTest.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, testToTest.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"},
			config.ExpressionImpl{Value: "link2"}}, testToTest.GetLinkerOptions().GetFlags())

		require.IsType(t, config.Dependencies{}, project.GetDependencies())
		require.IsType(t, &config.DependencyImpl{}, project.GetDependencies()["somedep"])
		require.Equal(t, &config.DependencyImpl{Ref: "default"}, project.GetDependencies()["somedep"])

		require.IsType(t, config.Dependencies{}, project.GetTestDependencies())
		require.IsType(t, &config.DependencyImpl{}, project.GetTestDependencies()["somedep"])
		require.Equal(t, &config.DependencyImpl{Ref: "default"}, project.GetTestDependencies()["somedep"])
	})

	suite.T().Run("Error - file not found", func(t *testing.T) {
		_, _, err := config.ReadConfig("randomFilePath")
		require.Error(t, err)
	})

	suite.T().Run("Error - invalid file content", func(t *testing.T) {
		_, _, err := config.ReadConfig(randomFileContentDir)
		require.Error(t, err)
	})

	suite.T().Run("Error - invalid schema", func(t *testing.T) {
		_, _, err := config.ReadConfig(invalidSchemaDir)
		require.Error(t, err)
	})

	suite.T().Run("Error - unsupported tag", func(t *testing.T) {
		_, _, err := config.ReadConfig(unsupportedTagDir)
		require.Error(t, err)
	})

	suite.T().Run("Happy path - hil usage", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(hilUsageDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		argument1, err := project.GetArguments()[0].GetValue(suite.config)

		variable1, err := project.GetVariables()[0].GetValue(suite.config)
		require.NoError(t, err)
		variable2, err := project.GetVariables()[1].GetValue(suite.config)
		require.NoError(t, err)
		variable3, err := project.GetVariables()[2].GetValue(suite.config)
		require.NoError(t, err)
		variable4, err := project.GetVariables()[3].GetValue(suite.config)
		require.NoError(t, err)
		variable5, err := project.GetVariables()[4].GetValue(suite.config)
		require.NoError(t, err)
		variable6, err := project.GetVariables()[5].GetValue(suite.config)
		require.NoError(t, err)
		variable7, err := project.GetVariables()[6].GetValue(suite.config)
		require.NoError(t, err)

		err = hillang.Initialize(project.GetVariables(), project.GetArguments())
		require.NoError(t, err)

		scripts := project.GetScripts()

		script1, err := scripts["begin"].Eval(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable3, script1)

		script2, err := scripts["end"].Eval(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable4, script2)

		projectName, err := project.GetProject().GetName(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable1, projectName)

		projectAuthor, err := project.GetProject().GetAuthor(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable2, projectAuthor)

		projectVersion, err := project.GetProject().GetVersion(suite.config)
		require.NoError(t, err)
		expectedVersion, err := version.NewVersion("0.0.1")
		require.NoError(t, err)
		require.Equal(t, expectedVersion, projectVersion)

		projectHomepage, err := project.GetProject().GetHomepage(suite.config)
		require.NoError(t, err)
		require.Equal(t, argument1, projectHomepage)

		projectDescription, err := project.GetProject().GetDescription(suite.config)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s project description", variable1), projectDescription)

		executableOptions := project.GetTargets()["main"].GetExecutableOptions()

		source1, err := executableOptions.GetSource()[0].Eval(suite.config)
		require.NoError(t, err)
		require.Equal(t, "src", source1)

		platform, err := executableOptions.GetPlatform(suite.config)
		require.NoError(t, err)
		require.Equal(t, "native", platform)

		targetScopeConfig, err := hillang.GetArgsEvalConfig(project.GetTargets()["main"].GetArguments(), suite.config)
		require.NoError(t, err)

		targetArgument1, err := project.GetTargets()["main"].GetArguments()[0].GetValue(suite.config)
		require.NoError(t, err)
		require.Equal(t, targetArgument1, "false")

		toolchainName, err := executableOptions.GetToolchain().GetName(targetScopeConfig)
		require.NoError(t, err)
		toolchainRef, err := executableOptions.GetToolchain().GetRef(targetScopeConfig)
		require.NoError(t, err)

		require.Equal(t, "prod", toolchainName)
		require.Equal(t, "default", toolchainRef)

		testTarget := project.GetTests()["main"]

		testScopeConfig, err := hillang.GetArgsEvalConfig(testTarget.GetArguments(), suite.config)
		require.Equal(t, config.Arguments{config.ArgumentImpl{Name: "VISIBILITY_CHECK", Value: "true"}},
			testTarget.GetArguments())

		require.NoError(t, err)

		targetName, err := testTarget.GetTargetName(testScopeConfig)
		require.NoError(t, err)
		require.Equal(t, "main", targetName)

		targetArgumentsOneName := testTarget.GetTargetArguments()[0].GetName()
		require.NoError(t, err)
		require.Equal(t, "DEBUG", targetArgumentsOneName)

		targetArgumentsOneValue, err := testTarget.GetTargetArguments()[0].GetValue(testScopeConfig)
		require.NoError(t, err)
		require.Equal(t, "true", targetArgumentsOneValue)

		linkerVisibility, err := testTarget.GetLinkerOptions().GetVisibility(testScopeConfig)
		require.NoError(t, err)
		require.Equal(t, variable5, linkerVisibility)

		dependency1 := project.GetDependencies()["dependency1"]
		dep1Ref, err := dependency1.GetRef(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable6, dep1Ref)

		dep1ArgumentOneName := dependency1.GetArguments()[0].GetName()
		require.Equal(t, "DEBUG", dep1ArgumentOneName)
		dep1ArgumentOneValue, err := dependency1.GetArguments()[0].GetValue(suite.config)
		require.NoError(t, err)
		require.Equal(t, "true", dep1ArgumentOneValue)

		dep1LinkerVisibility, err := dependency1.GetLinkerOptions().GetVisibility(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable7, dep1LinkerVisibility)

		testDependency1 := project.GetDependencies()["dependency1"]
		testDep1Ref, err := testDependency1.GetRef(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable6, testDep1Ref)

		testDep1ArgumentOneName := testDependency1.GetArguments()[0].GetName()
		require.Equal(t, "DEBUG", testDep1ArgumentOneName)
		testDep1ArgumentOneValue, err := testDependency1.GetArguments()[0].GetValue(suite.config)
		require.NoError(t, err)
		require.Equal(t, "true", testDep1ArgumentOneValue)

		testDep1LinkerVisibility, err := testDependency1.GetLinkerOptions().GetVisibility(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable7, testDep1LinkerVisibility)

		require.Equal(t, nil, project.GetTargets()["main"].GetCompileOptions())
		require.Equal(t, nil, project.GetProject().GetCompileOptions())
	})

	suite.T().Run("Error - invalid hil usage and invalid version", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(invalidHilUsageDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		err = hillang.Initialize(project.GetVariables(), project.GetArguments())
		require.NoError(t, err)

		// invalid function
		_, err = project.GetScripts()["begin"].Eval(suite.config)
		require.Error(t, err)

		// invalid syntax
		_, err = project.GetProject().GetName(suite.config)
		require.Error(t, err)

		// invalid version
		_, err = project.GetProject().GetVersion(suite.config)
		require.Error(t, err)

		// invalid boolean
		_, err = project.GetProject().GetPackageOptions().IsHeaderOnly(suite.config)
		require.Error(t, err)
	})

	suite.T().Run("Happy path - script exec for fields", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(configScriptExecDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		variable1, err := project.GetVariables()[0].GetValue(suite.config)
		require.NoError(t, err)

		err = hillang.Initialize(project.GetVariables(), project.GetArguments())
		require.NoError(t, err)

		projectName, err := project.GetProject().GetName(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable1, projectName)

		projectDescription, err := project.GetProject().GetDescription(suite.config)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s project description", variable1), projectDescription)

		executableOptions := project.GetTargets()["main"].GetExecutableOptions()

		source1, err := executableOptions.GetSource()[0].Eval(suite.config)
		require.NoError(t, err)
		require.Equal(t, "src", source1)

		platform, err := executableOptions.GetPlatform(suite.config)
		require.NoError(t, err)
		require.Equal(t, "native", platform)
	})

	suite.T().Run("Error - script exec invalid", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(configScriptExecInvalidDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		err = hillang.Initialize(project.GetVariables(), project.GetArguments())
		require.NoError(t, err)

		// invalid exec function
		_, err = project.GetProject().GetName(suite.config)
		require.Error(t, err)

		// invalid version variable
		_, err = project.GetProject().GetVersion(suite.config)
		require.Error(t, err)

		// invalid hil eval
		_, err = project.GetProject().GetDescription(suite.config)
		require.Error(t, err)

		executableOptions := project.GetTests()["main"].GetExecutableOptions()

		// invalid variable
		_, err = executableOptions.GetSource()[0].Eval(suite.config)
		require.Error(t, err)

		// invalid boolean
		_, err = project.GetProject().GetPackageOptions().IsHeaderOnly(suite.config)
		require.Error(t, err)
	})
}

func (suite *ConfigTestSuite) TestCreateConfig() {
	// happy path app
	suite.T().Run("Happy path - create app config with toolchain", func(t *testing.T) {
		creationConfig := templates.ProjectCreation{
			Type:        constants.APP,
			ProjectName: "AppWithToolchain",
			ProjectPath: "/AppWithToolchain",
			Platform:    "native",
			Toolchain:   "clang",
			MainFile:    "src/main.cpp",
		}

		err := config.CreateConfig(creationConfig)
		require.NoError(t, err)

		createdContent, err := sys.ReadFile(fmt.Sprintf("%s/%s",
			creationConfig.ProjectPath, constants.WioConfigFile))
		require.NoError(t, err)

		require.Equal(t,
			fmt.Sprintf(createConfigAppWithToolchain, creationConfig.ProjectName, creationConfig.MainFile,
				creationConfig.Platform, creationConfig.Toolchain, creationConfig.Platform, creationConfig.Toolchain),
			string(createdContent))

		parsedConfig, warnings, err := config.ReadConfig(creationConfig.ProjectPath)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, constants.APP, parsedConfig.GetType())
	})

	suite.T().Run("Happy path - create app config without toolchain", func(t *testing.T) {
		creationConfig := templates.ProjectCreation{
			Type:        constants.APP,
			ProjectName: "AppWithoutToolchain",
			ProjectPath: "/AppWithoutToolchain",
			Platform:    "native",
			MainFile:    "src/main.cpp",
		}

		err := config.CreateConfig(creationConfig)
		require.NoError(t, err)

		createdContent, err := sys.ReadFile(fmt.Sprintf("%s/%s",
			creationConfig.ProjectPath, constants.WioConfigFile))
		require.NoError(t, err)

		require.Equal(t,
			fmt.Sprintf(createConfigAppWithoutToolchain, creationConfig.ProjectName, creationConfig.MainFile,
				creationConfig.Platform, creationConfig.Platform),
			string(createdContent))

		parsedConfig, warnings, err := config.ReadConfig(creationConfig.ProjectPath)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, constants.APP, parsedConfig.GetType())
	})

	// happy path pkg
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
