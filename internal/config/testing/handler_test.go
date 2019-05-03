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

		target := project.GetTargets()[0]
		targetName, err := target.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		require.Equal(t, project.GetProject().GetPackageOptions(), target.GetPackageOptions())
		require.Equal(t, project.GetProject().GetCompileOptions(), target.GetCompileOptions())
	})

	suite.T().Run("Happy path - some scope values provided, should add others from global", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(someScopeValuesDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		targetMain2 := project.GetTargets()[1]
		targetName, err = targetMain2.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main2", targetName)

		expectedb, err := project.GetProject().GetPackageOptions().IsHeaderOnly(suite.config)
		require.NoError(t, err)
		actualb, err := targetMain.GetPackageOptions().IsHeaderOnly(suite.config)
		require.NoError(t, err)

		require.Equal(t, expectedb, actualb)

		actualb, err = targetMain2.GetPackageOptions().IsHeaderOnly(suite.config)
		require.NoError(t, err)

		require.Equal(t, false, actualb)

		actual, err := targetMain.GetPackageOptions().GetPackageType(suite.config)
		require.NoError(t, err)
		require.Equal(t, "SHARED", actual)

		require.Equal(t, project.GetProject().GetCompileOptions().GetFlags(),
			targetMain.GetCompileOptions().GetFlags())
		require.Equal(t, project.GetProject().GetCompileOptions().GetDefinitions(),
			targetMain.GetCompileOptions().GetDefinitions())

		actual, err = targetMain.GetCompileOptions().GetCXXStandard(suite.config)
		require.NoError(t, err)
		require.Equal(t, "c++17", actual)

		actual, err = targetMain.GetCompileOptions().GetCStandard(suite.config)
		require.NoError(t, err)
		require.Equal(t, "c01", actual)
	})

	suite.T().Run("Happy path - array values are provided, append them with global", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(arrayValuesProvidedDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		globalFlags := project.GetProject().GetCompileOptions().GetFlags()
		targetFlags := targetMain.GetCompileOptions().GetFlags()

		require.Equal(t, append(config.Flags{config.ExpressionImpl{Value: "flag3"}}, globalFlags...), targetFlags)

		globalDefinitions := project.GetProject().GetCompileOptions().GetDefinitions()
		targetDefinitions := targetMain.GetCompileOptions().GetDefinitions()

		require.Equal(t, append(config.Definitions{config.ExpressionImpl{Value: "def3"}}, globalDefinitions...),
			targetDefinitions)

		globalCXX, err := project.GetProject().GetCompileOptions().GetCXXStandard(suite.config)
		require.NoError(t, err)
		globalC, err := project.GetProject().GetCompileOptions().GetCStandard(suite.config)
		require.NoError(t, err)

		targetCXX, err := targetMain.GetCompileOptions().GetCXXStandard(suite.config)
		require.NoError(t, err)
		targetC, err := targetMain.GetCompileOptions().GetCStandard(suite.config)
		require.NoError(t, err)

		require.Equal(t, globalCXX, targetCXX)
		require.Equal(t, globalC, targetC)
	})

	suite.T().Run("Warnings - app config warnings", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(appConfigWarningsDir)
		require.NoError(t, err)

		require.Equal(t, "app", project.GetType())
		require.Equal(t, 3, len(warnings))

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		testMain := project.GetTests()[0]
		testName, err := testMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", testName)

		require.Equal(t, nil, project.GetProject().GetPackageOptions())
		require.Equal(t, nil, targetMain.GetPackageOptions())

		file, err := testMain.GetExecutableOptions().GetMainFile(suite.config)
		require.NoError(t, err)

		require.Equal(t, "", file)
	})

	suite.T().Run("Warnings - pkg config warnings", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(pkgConfigWarningsDir)
		require.NoError(t, err)

		require.Equal(t, "pkg", project.GetType())
		require.Equal(t, 2, len(warnings))

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		testMain := project.GetTests()[0]
		testName, err := testMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", testName)

		require.Equal(t, nil, targetMain.GetExecutableOptions())

		file, err := testMain.GetExecutableOptions().GetMainFile(suite.config)
		require.NoError(t, err)

		require.Equal(t, "", file)
	})

	suite.T().Run("Happy path - convert a string to string array for certain fields", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(stringToSliceFieldsDir)
		require.NoError(t, err)

		require.Equal(t, 0, len(warnings))

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		targetMain2 := project.GetTargets()[1]
		targetName, err = targetMain2.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "bru", targetName)

		testMain := project.GetTests()[0]
		testName, err := testMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", testName)

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

		require.IsType(t, config.Sources{}, targetMain.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "src"}},
			targetMain.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, targetMain.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "github.com/wio-pm/someToolchain", Ref: "develop"},
			targetMain.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, targetMain.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}, config.ExpressionImpl{Value: "flag1"}},
			targetMain.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, targetMain.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}, config.ExpressionImpl{Value: "def1"}},
			targetMain.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, targetMain.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"}}, targetMain.GetLinkerOptions().GetFlags())

		require.IsType(t, config.Sources{}, testMain.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "test"}},
			testMain.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, testMain.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "github.com/wio-pm/someToolchain", Ref: ""},
			testMain.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, testMain.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}},
			testMain.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, testMain.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}},
			testMain.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, testMain.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"}}, testMain.GetLinkerOptions().GetFlags())

		dependency := project.GetDependencies()[0]
		dependencyName, err := dependency.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "github.com/random/dep", dependencyName)

		testDependency := project.GetTestDependencies()[0]
		testDependencyName, err := testDependency.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "github.com/random/test_dep", testDependencyName)

		require.IsType(t, config.Dependencies{}, project.GetDependencies())
		require.Equal(t, &config.DependencyImpl{Name: dependencyName, Ref: ""}, dependency)

		require.IsType(t, config.Dependencies{}, project.GetTestDependencies())
		require.Equal(t, &config.DependencyImpl{Name: testDependencyName, Ref: "develop"}, testDependency)
	})

	suite.T().Run("Happy path - convert a string to string array if separated by ,", func(t *testing.T) {
		project, warnings, err := config.ReadConfig(stringToSliceFieldsCommasDir)
		require.NoError(t, err)

		require.Equal(t, 0, len(warnings))

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		testMain := project.GetTests()[0]
		testName, err := testMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", testName)

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

		require.IsType(t, config.Sources{}, targetMain.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "src"},
			config.ExpressionImpl{Value: "common"}, config.ExpressionImpl{Value: "utils"}},
			targetMain.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, targetMain.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "someToolchain", Ref: "default"},
			targetMain.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, targetMain.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}, config.ExpressionImpl{Value: "flag4"},
			config.ExpressionImpl{Value: "flag1"}, config.ExpressionImpl{Value: "flag2"}},
			targetMain.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, targetMain.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}, config.ExpressionImpl{Value: "def4"},
			config.ExpressionImpl{Value: "def1"}, config.ExpressionImpl{Value: "def2"}},
			targetMain.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, targetMain.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"},
			config.ExpressionImpl{Value: "link2"}}, targetMain.GetLinkerOptions().GetFlags())

		require.IsType(t, config.Sources{}, testMain.GetExecutableOptions().GetSource())
		require.Equal(t, config.Sources{config.ExpressionImpl{Value: "test"}, config.ExpressionImpl{Value: "utils"}},
			testMain.GetExecutableOptions().GetSource())

		require.IsType(t, config.ToolchainImpl{}, testMain.GetExecutableOptions().GetToolchain())
		require.Equal(t, config.ToolchainImpl{Name: "someToolchain", Ref: "default"},
			testMain.GetExecutableOptions().GetToolchain())

		require.IsType(t, config.Flags{}, testMain.GetCompileOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "flag2"}, config.ExpressionImpl{Value: "flag4"}},
			testMain.GetCompileOptions().GetFlags())

		require.IsType(t, config.Definitions{}, testMain.GetCompileOptions().GetDefinitions())
		require.Equal(t, config.Definitions{config.ExpressionImpl{Value: "def2"}, config.ExpressionImpl{Value: "def4"}},
			testMain.GetCompileOptions().GetDefinitions())

		require.IsType(t, config.Flags{}, testMain.GetLinkerOptions().GetFlags())
		require.Equal(t, config.Flags{config.ExpressionImpl{Value: "link1"},
			config.ExpressionImpl{Value: "link2"}}, testMain.GetLinkerOptions().GetFlags())
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

		targetMain := project.GetTargets()[0]
		targetName, err := targetMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", targetName)

		testMain := project.GetTests()[0]
		testName, err := testMain.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "main", testName)

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

		executableOptions := targetMain.GetExecutableOptions()

		source1, err := executableOptions.GetSource()[0].Eval(suite.config)
		require.NoError(t, err)
		require.Equal(t, "src", source1)

		platform, err := executableOptions.GetPlatform(suite.config)
		require.NoError(t, err)
		require.Equal(t, "native", platform)

		targetScopeConfig, err := hillang.GetArgsEvalConfig(targetMain.GetArguments(), suite.config)
		require.NoError(t, err)

		targetArgument1, err := targetMain.GetArguments()[0].GetValue(suite.config)
		require.NoError(t, err)
		require.Equal(t, targetArgument1, "false")

		toolchainName, err := executableOptions.GetToolchain().GetName(targetScopeConfig)
		require.NoError(t, err)
		toolchainRef, err := executableOptions.GetToolchain().GetRef(targetScopeConfig)
		require.NoError(t, err)

		require.Equal(t, "prod", toolchainName)
		require.Equal(t, "default", toolchainRef)

		testScopeConfig, err := hillang.GetArgsEvalConfig(testMain.GetArguments(), suite.config)
		require.Equal(t, config.Arguments{config.ArgumentImpl{Name: "VISIBILITY_CHECK", Value: "true"}},
			testMain.GetArguments())

		require.NoError(t, err)

		testTargetName, err := testMain.GetTargetName(testScopeConfig)
		require.NoError(t, err)
		require.Equal(t, "main", testTargetName)

		targetArgumentsOneName := testMain.GetTargetArguments()[0].GetName()
		require.NoError(t, err)
		require.Equal(t, "DEBUG", targetArgumentsOneName)

		targetArgumentsOneValue, err := testMain.GetTargetArguments()[0].GetValue(testScopeConfig)
		require.NoError(t, err)
		require.Equal(t, "true", targetArgumentsOneValue)

		linkerVisibility, err := testMain.GetLinkerOptions().GetVisibility(testScopeConfig)
		require.NoError(t, err)
		require.Equal(t, variable5, linkerVisibility)

		dependency := project.GetDependencies()[0]
		dependencyName, err := dependency.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "gitlab.com/user/dependency1", dependencyName)

		testDependency := project.GetTestDependencies()[0]
		testDependencyName, err := testDependency.GetName(suite.config)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), "gitlab.com/user/dependency1", testDependencyName)

		dep1Ref, err := dependency.GetRef(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable6, dep1Ref)

		dep1ArgumentOneName := dependency.GetArguments()[0].GetName()
		require.Equal(t, "DEBUG", dep1ArgumentOneName)
		dep1ArgumentOneValue, err := dependency.GetArguments()[0].GetValue(suite.config)
		require.NoError(t, err)
		require.Equal(t, "true", dep1ArgumentOneValue)

		dep1LinkerVisibility, err := dependency.GetLinkerOptions().GetVisibility(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable7, dep1LinkerVisibility)

		testDep1Ref, err := testDependency.GetRef(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable6, testDep1Ref)

		testDep1ArgumentOneName := testDependency.GetArguments()[0].GetName()
		require.Equal(t, "DEBUG", testDep1ArgumentOneName)
		testDep1ArgumentOneValue, err := testDependency.GetArguments()[0].GetValue(suite.config)
		require.NoError(t, err)
		require.Equal(t, "true", testDep1ArgumentOneValue)

		testDep1LinkerVisibility, err := testDependency.GetLinkerOptions().GetVisibility(suite.config)
		require.NoError(t, err)
		require.Equal(t, variable7, testDep1LinkerVisibility)

		// they just needed to be tested
		require.Equal(t, nil, targetMain.GetCompileOptions())
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

		executableOptions := project.GetTargets()[0].GetExecutableOptions()

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

		executableOptions := project.GetTests()[0].GetExecutableOptions()

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

	suite.T().Run("Happy path - create pkg config header only with toolchain", func(t *testing.T) {
		creationConfig := templates.ProjectCreation{
			Type:        constants.PKG,
			ProjectName: "PkgHeaderOnlyWithToolchain",
			ProjectPath: "/PkgHeaderOnlyWithToolchain",
			Platform:    "native",
			Toolchain:   "clang",
			HeaderOnly:  true,
		}

		err := config.CreateConfig(creationConfig)
		require.NoError(t, err)

		createdContent, err := sys.ReadFile(fmt.Sprintf("%s/%s",
			creationConfig.ProjectPath, constants.WioConfigFile))
		require.NoError(t, err)

		require.Equal(t,
			fmt.Sprintf(createConfigPkgHeaderOnlyToolchain, creationConfig.ProjectName, creationConfig.ProjectName,
				creationConfig.ProjectName, creationConfig.Platform,
				creationConfig.Toolchain, creationConfig.ProjectName),
			string(createdContent))

		parsedConfig, warnings, err := config.ReadConfig(creationConfig.ProjectPath)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, constants.PKG, parsedConfig.GetType())
	})

	suite.T().Run("Happy path - create pkg config without toolchain", func(t *testing.T) {
		creationConfig := templates.ProjectCreation{
			Type:        constants.PKG,
			ProjectName: "PkgWithoutToolchain",
			ProjectPath: "/PkgWithoutToolchain",
			Platform:    "native",
		}

		err := config.CreateConfig(creationConfig)
		require.NoError(t, err)

		createdContent, err := sys.ReadFile(fmt.Sprintf("%s/%s",
			creationConfig.ProjectPath, constants.WioConfigFile))
		require.NoError(t, err)

		require.Equal(t,
			fmt.Sprintf(createConfigPkgWithoutToolchain, creationConfig.ProjectName, creationConfig.ProjectName,
				creationConfig.ProjectName, creationConfig.Platform, creationConfig.ProjectName),
			string(createdContent))

		parsedConfig, warnings, err := config.ReadConfig(creationConfig.ProjectPath)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, constants.PKG, parsedConfig.GetType())
	})

	suite.T().Run("Happy path - create pkg config header only shared", func(t *testing.T) {
		creationConfig := templates.ProjectCreation{
			Type:        constants.PKG,
			ProjectName: "PkgHeaderOnlyShared",
			ProjectPath: "/PkgHeaderOnlyShared",
			Platform:    "native",
			HeaderOnly:  true,
			PkgType:     "SHARED",
		}

		err := config.CreateConfig(creationConfig)
		require.NoError(t, err)

		createdContent, err := sys.ReadFile(fmt.Sprintf("%s/%s",
			creationConfig.ProjectPath, constants.WioConfigFile))
		require.NoError(t, err)

		require.Equal(t,
			fmt.Sprintf(createConfigPkgHeaderOnlyShared, creationConfig.ProjectName, creationConfig.ProjectName,
				creationConfig.ProjectName, creationConfig.Platform, creationConfig.ProjectName),
			string(createdContent))

		parsedConfig, warnings, err := config.ReadConfig(creationConfig.ProjectPath)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, constants.PKG, parsedConfig.GetType())
	})

	suite.T().Run("Happy path - create pkg config shared", func(t *testing.T) {
		creationConfig := templates.ProjectCreation{
			Type:        constants.PKG,
			ProjectName: "PkgShared",
			ProjectPath: "/PkgShared",
			Platform:    "native",
			PkgType:     "SHARED",
		}

		err := config.CreateConfig(creationConfig)
		require.NoError(t, err)

		createdContent, err := sys.ReadFile(fmt.Sprintf("%s/%s",
			creationConfig.ProjectPath, constants.WioConfigFile))
		require.NoError(t, err)

		require.Equal(t,
			fmt.Sprintf(createConfigPkgShared, creationConfig.ProjectName, creationConfig.ProjectName,
				creationConfig.ProjectName, creationConfig.Platform, creationConfig.ProjectName),
			string(createdContent))

		parsedConfig, warnings, err := config.ReadConfig(creationConfig.ProjectPath)
		require.NoError(t, err)
		require.Equal(t, 0, len(warnings))

		require.Equal(t, constants.PKG, parsedConfig.GetType())
	})

	// happy path pkg
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
