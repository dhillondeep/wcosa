package hillang

import (
	"github.com/hashicorp/hil/ast"
	"github.com/huandu/xstrings"
	"github.com/thoas/go-funk"
	"os"
	"strconv"
	"strings"
)

var env = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return os.Getenv(input), nil
	},
}

var lowerCase = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return strings.ToLower(input), nil
	},
}

var upperCase = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return strings.ToUpper(input), nil
	},
}

var snakeCase = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return xstrings.ToSnakeCase(input), nil
	},
}

var camelCase = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return xstrings.ToCamelCase(input), nil
	},
}

var reverse = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return xstrings.Reverse(input), nil
	},
}

var shuffle = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return xstrings.Shuffle(input), nil
	},
}

var wordCount = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeInt,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return xstrings.WordCount(input), nil
	},
}

var length = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeInt,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(string)
		return xstrings.Len(input), nil
	},
}

var toString = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeInt},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		input := inputs[0].(int)
		return strconv.Itoa(input), nil
	},
}

var join = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeList, ast.TypeString},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		list := inputs[0].([]string)
		sep := inputs[1].(string)

		return strings.Join(list, sep), nil
	},
}

var insert = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString, ast.TypeString, ast.TypeInt},
	ReturnType: ast.TypeString,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		dst := inputs[0].(string)
		src := inputs[1].(string)
		index := inputs[2].(int)
		return xstrings.Insert(dst, src, index), nil
	},
}

var defined = ast.Function{
	ArgTypes:   []ast.Type{ast.TypeString},
	ReturnType: ast.TypeBool,
	Variadic:   false,
	Callback: func(inputs []interface{}) (interface{}, error) {
		originalInput := inputs[0].(string)
		split := strings.Split(originalInput, ".")

		if len(split) < 2 {
			return false, nil
		}

		if split[0] == "var" {
			return funk.Contains(variablesMap, originalInput), nil
		} else if split[0] == "arg" {
			return funk.Contains(argsMap, originalInput), nil
		}

		return false, nil
	},
}

func getFunctions() map[string]ast.Function {
	return map[string]ast.Function{
		"env": env,
		"lower":  lowerCase,
		"upper":  upperCase,
		"snakeCase": snakeCase,
		"camelCase": camelCase,
		"reverse": reverse,
		"shuffle": shuffle,
		"wordCount": wordCount,
		"length": length,
		"string": toString,
		"join": join,
		"insert": insert,
		"defined": defined,
	}
}
