package tengolang

import (
	"fmt"
	"github.com/d5/tengo/script"
	"github.com/huandu/xstrings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getScript(t *testing.T, scriptData string) *script.Script {
	s := script.New([]byte(scriptData))
	s.SetImports(GetModuleMap("wstrings"))
	if err := s.Add("out", ""); err != nil {
		assert.Fail(t, err.Error())
	}

	return s
}

func TestGetModuleMap(t *testing.T) {
	// all Modules
	allModules := fmt.Sprintf(`
	math := import("math")
	wstrings := import("wstrings")

	out = wstrings.snakeCase(string(math.abs(%d)) + " Brother")
	`, -10)

	s := script.New([]byte(allModules))
	s.SetImports(GetModuleMap("math", "wstrings"))
	if err := s.Add("out", ""); err != nil {
		assert.Fail(t, err.Error())
	}

	c, err := s.Run()

	assert.NoError(t, err)
	assert.Equal(t, "10_brother", c.Get("out").String())

	// stdlib Modules
	stdlibModules := fmt.Sprintf(`
	math := import("math")

	out = math.ceil(%f)
	`, 11.98)

	s = script.New([]byte(stdlibModules))
	s.SetImports(GetModuleMap("math"))
	if err := s.Add("out", ""); err != nil {
		assert.Fail(t, err.Error())
	}

	c, err = s.Run()

	assert.NoError(t, err)
	assert.Equal(t, 12, c.Get("out").Int())
}

func TestSnakeCaseFunc(t *testing.T) {
	word := "Value"

	// happy path
	tengoScript := fmt.Sprintf(`
	wstrings := import("wstrings")

	out = wstrings.snakeCase("%s")
	`, word)

	s := getScript(t, tengoScript)
	c, err := s.Run()

	assert.NoError(t, err)
	assert.Equal(t, xstrings.ToSnakeCase(word), c.Get("out").String())

	// wrong input - undefined type
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.snakeCase(func() { b := 4 }())
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)

	// wrong input - wrong number of arguments
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.snakeCase()
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)
}

func TestCamelCaseFunc(t *testing.T) {
	// happy path
	word := "Value"
	tengoScript := fmt.Sprintf(`
	wstrings := import("wstrings")

	out = wstrings.camelCase("%s")
	`, word)

	s := getScript(t, tengoScript)
	c, err := s.Run()

	assert.NoError(t, err)
	assert.Equal(t, xstrings.ToCamelCase(word), c.Get("out").String())

	// wrong input - undefined type
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.camelCase(func() { b := 4 }())
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)

	// wrong input - wrong number of arguments
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.camelCase()
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)
}

func TestReverseCaseFunc(t *testing.T) {
	// happy path
	word := "Value"
	tengoScript := fmt.Sprintf(`
	wstrings := import("wstrings")

	out = wstrings.reverse("%s")
	`, word)

	s := getScript(t, tengoScript)
	c, err := s.Run()

	assert.NoError(t, err)
	assert.Equal(t, xstrings.Reverse(word), c.Get("out").String())

	// wrong input - undefined type
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.reverse(func() { b := 4 }())
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)

	// wrong input - wrong number of arguments
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.reverse()
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)
}

func TestShuffleCaseFunc(t *testing.T) {
	// happy path
	word := "Value"
	tengoScript := fmt.Sprintf(`
	wstrings := import("wstrings")

	out = wstrings.shuffle("%s")
	`, word)

	s := getScript(t, tengoScript)
	c, err := s.Run()

	assert.NoError(t, err)
	assert.NotEqual(t, word, c.Get("out").String())

	// wrong input - undefined type
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.shuffle(func() { b := 4 }())
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)

	// wrong input - wrong number of arguments
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.shuffle()
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)
}

func TestWordCountFunc(t *testing.T) {
	// happy path
	word := "Value of Bro"
	tengoScript := fmt.Sprintf(`
	wstrings := import("wstrings")

	out = wstrings.wordCount("%s")
	`, word)

	s := getScript(t, tengoScript)
	c, err := s.Run()

	assert.NoError(t, err)
	assert.Equal(t, xstrings.WordCount(word), c.Get("out").Int())

	// wrong input - undefined type
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.wordCount(func() { b := 4 }())
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)

	// wrong input - wrong number of arguments
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.wordCount()
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)
}

func TestLengthFunc(t *testing.T) {
	// happy path
	word := "Value of Bro"
	tengoScript := fmt.Sprintf(`
	wstrings := import("wstrings")

	out = wstrings.length("%s")
	`, word)

	s := getScript(t, tengoScript)
	c, err := s.Run()

	assert.NoError(t, err)
	assert.Equal(t, xstrings.Len(word), c.Get("out").Int())

	// wrong input - undefined type
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.length(func() { b := 4 }())
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)

	// wrong input - wrong number of arguments
	tengoScript = `
	wstrings := import("wstrings")

	out = wstrings.length()
	`

	s = getScript(t, tengoScript)
	c, err = s.Run()

	assert.Error(t, err)
}
