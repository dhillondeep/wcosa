package sys

import (
	"github.com/dhillondeep/afero"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestWriteJson(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	t.Run("Happy path - valid Json", func(t *testing.T) {
		type sample struct {
			Hello string
		}

		expected := `{
  "Hello": "World"
}`

		err := WriteJson("someRandomFile.json", &sample{Hello: "World"})
		require.NoError(t, err)

		fileData, err := ReadFile("someRandomFile.json")
		require.NoError(t, err)

		require.Equal(t, expected, string(fileData))
	})

	t.Run("Error - invalid json", func(t *testing.T) {
		err := WriteJson("someRandomFile.json", math.NaN())
		require.Error(t, err)
	})
}

func TestWriteYaml(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	t.Run("Happy path - valid Json", func(t *testing.T) {
		type sample struct {
			Hello string
		}

		expected := `hello: World
`

		err := WriteYaml("someRandomFile.yaml", &sample{Hello: "World"})
		require.NoError(t, err)

		fileData, err := ReadFile("someRandomFile.yaml")
		require.NoError(t, err)

		require.Equal(t, expected, string(fileData))
	})

	t.Run("Error - invalid yaml panic", func(t *testing.T) {
		type sample struct {
			Hello string `yaml:"hello"`
			Bye   string `yaml:"hello"`
		}

		require.Panics(t, func() {
			_ = WriteYaml("someRandomFile.yaml", &sample{Hello: "world", Bye: "world"})
		})
	})

	t.Run("Error - invalid yaml error", func(t *testing.T) {
		type sample struct {
			Hello string `yaml:"hello"`
			Bye   string `yaml:"hello"`
		}

		oldYamlMarshal := yamlMarshal
		defer func() {
			yamlMarshal = oldYamlMarshal
		}()
		yamlMarshal = func(in interface{}) (out []byte, err error) {
			return nil, errors.New("some error")
		}

		err := WriteYaml("someRandomFile.yaml", &sample{Hello: "world", Bye: "world"})
		require.Error(t, err)
	})
}

func TestParseJson(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	t.Run("Happy path - successful parse", func(t *testing.T) {
		type sampleJson struct {
			Hello string
		}

		data := []byte(`{
  "Hello": "World"
}`)

		require.NoError(t, WriteFile("someJsonFile.json", data))

		sampleJsonData := &sampleJson{}
		require.NoError(t, ParseJson("someJsonFile.json", sampleJsonData))
		require.Equal(t, &sampleJson{
			Hello: "World",
		}, sampleJsonData)
	})

	t.Run("Happy path - unsuccessful parse", func(t *testing.T) {
		type sample struct {
			Hello string
		}

		data := []byte(`{
  "no": "world"
}`)

		require.NoError(t, WriteFile("someJsonFile.json", data))

		sampleData := &sample{}
		require.NoError(t, ParseJson("someJsonFile.json", sampleData))
		require.Equal(t, &sample{
			Hello: "",
		}, sampleData)
	})

	t.Run("Error - file does not exist", func(t *testing.T) {
		require.Error(t, ParseJson("someJsonFileRandom.json", map[string]string{}))
	})

	t.Run("Error - not valid", func(t *testing.T) {
		data := []byte(`random,bru,{"Hello"}`)

		require.NoError(t, WriteFile("someJsonFile.json", data))

		require.Error(t, ParseJson("someJsonFile.json", map[string]string{}))
	})
}

func TestParseYaml(t *testing.T) {
	SetFileSystem(afero.NewMemMapFs())

	t.Run("Happy path - successful parse", func(t *testing.T) {
		type sample struct {
			Hello string
		}

		data := []byte(`hello: World
`)

		require.NoError(t, WriteFile("someYamlFile.yaml", data))

		sampleData := &sample{}
		require.NoError(t, ParseYaml("someYamlFile.yaml", sampleData))
		require.Equal(t, &sample{
			Hello: "World",
		}, sampleData)
	})

	t.Run("Happy path - unsuccessful parse", func(t *testing.T) {
		type sample struct {
			Hello string
		}

		data := []byte(`no: World
`)

		require.NoError(t, WriteFile("someYamlFile.yaml", data))

		sampleData := &sample{}
		require.NoError(t, ParseYaml("someYamlFile.yaml", sampleData))
		require.Equal(t, &sample{
			Hello: "",
		}, sampleData)
	})

	t.Run("Error - file does not exist", func(t *testing.T) {
		require.Error(t, ParseYaml("someYamlFileRandom.Yaml", map[string]string{}))
	})

	t.Run("Error - not valid", func(t *testing.T) {
		data := []byte(`random,bru,{"Hello"}`)

		require.NoError(t, WriteFile("someYamlFile.yaml", data))

		require.Error(t, ParseYaml("someYamlFile.yaml", map[string]string{}))
	})
}
