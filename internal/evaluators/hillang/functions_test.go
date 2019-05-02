package hillang

import (
	"fmt"
	"github.com/hashicorp/hil"
	"github.com/huandu/xstrings"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
	"wio/internal/config"
)

func parseHil(val string, evalConfig *hil.EvalConfig) (string, error) {
	tree, err := hil.Parse(val)
	if err != nil {
		return "", err
	}

	result, err := hil.Eval(tree, evalConfig)
	if err != nil {
		return "", err
	}

	return result.Value.(string), nil
}

func TestEnv(t *testing.T) {
	if err := os.Setenv("JUST_DEFINED", "Value"); err != nil {
		assert.Fail(t, err.Error())
	}

	hilStatement := `${env("JUST_DEFINED")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.Equal(t, os.Getenv("JUST_DEFINED"), result)
}

func TestLowerCase(t *testing.T) {
	word := "Value"
	hilStatement := `${lower("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.Equal(t, strings.ToLower(word), result)
}

func TestUpperCase(t *testing.T) {
	word := "Value"
	hilStatement := `${upper("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.Equal(t, strings.ToUpper(word), result)
}

func TestSnakeCase(t *testing.T) {
	word := "Value For"
	hilStatement := `${snakeCase("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.Equal(t, xstrings.ToSnakeCase(word), result)
}

func TestCamelCase(t *testing.T) {
	word := "Value For"
	hilStatement := `${camelCase("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.Equal(t, xstrings.ToCamelCase(word), result)
}

func TestReverse(t *testing.T) {
	word := "Value For"
	hilStatement := `${reverse("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.Equal(t, xstrings.Reverse(word), result)
}

func TestShuffle(t *testing.T) {
	word := "Value For"
	hilStatement := `${shuffle("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)
	assert.NotEqual(t, word, result)
}

func TestWordCount(t *testing.T) {
	word := "Value For"
	hilStatement := `${wordCount("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)

	num, err := strconv.Atoi(result)
	assert.NoError(t, err)

	assert.Equal(t, xstrings.WordCount(word), num)
}

func TestLength(t *testing.T) {
	word := "Value For"
	hilStatement := `${length("` + word + `")}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)

	num, err := strconv.Atoi(result)
	assert.NoError(t, err)

	assert.Equal(t, xstrings.Len(word), num)
}

func TestToString(t *testing.T) {
	word := "Value"
	hilStatement := `${toString(length("` + word + `"))}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)

	assert.Equal(t, strconv.Itoa(xstrings.Len(word)), result)
}

func TestAppend(t *testing.T) {
	word := "Value"
	hilStatement := fmt.Sprintf(`${append("%s","%s", "%s")}`, word, word, word)

	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)

	assert.Equal(t, strings.Join(append([]string{}, word, word, word), ""), result)
}

func TestInsert(t *testing.T) {
	word := "Value"
	hilStatement := `${insert("` + word + `", "` + word + `", 2)}`
	result, err := parseHil(hilStatement, evalConfig)
	assert.NoError(t, err)

	assert.Equal(t, xstrings.Insert(word, word, 2), result)
}

func TestDefined(t *testing.T) {
	t.Run("happy path - none defined for variables and arguments", func(t *testing.T) {
		hilStatement := `${defined("var.RANDOM")}`
		result, err := parseHil(hilStatement, evalConfig)
		assert.NoError(t, err)

		assert.Equal(t, result, "false")

		hilStatement = `${defined("arg.RANDOM")}`
		result, err = parseHil(hilStatement, evalConfig)
		assert.NoError(t, err)

		assert.Equal(t, result, "false")
	})

	t.Run("happy path - variables and arguments defined", func(t *testing.T) {
		_ = Initialize(config.Variables{
			config.VariableImpl{
				Name:  "VARIABLE",
				Value: "One",
			},
		}, config.Arguments{
			config.ArgumentImpl{
				Name:  "ARGUMENT",
				Value: "One",
			},
		})

		hilStatement := `${defined("var.VARIABLE")}`
		result, err := parseHil(hilStatement, evalConfig)
		assert.NoError(t, err)

		assert.Equal(t, result, "true")

		hilStatement = `${defined("arg.ARGUMENT")}`
		result, err = parseHil(hilStatement, evalConfig)
		assert.NoError(t, err)

		assert.Equal(t, result, "true")
	})

	t.Run("wrong path - variable and argument scope not provided", func(t *testing.T) {
		hilStatement := `${defined("VARIABLE")}`
		result, err := parseHil(hilStatement, evalConfig)
		assert.NoError(t, err)
		assert.Equal(t, result, "false")
	})

	t.Run("wrong path - some random scope is defined", func(t *testing.T) {
		hilStatement := `${defined("test.VARIABLE")}`
		result, err := parseHil(hilStatement, evalConfig)
		assert.NoError(t, err)
		assert.Equal(t, result, "false")
	})
}
