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

var appendFunc = ast.Function{
	ArgTypes:     []ast.Type{ast.TypeString},
	ReturnType:   ast.TypeString,
	Variadic:     true,
	VariadicType: ast.TypeString,
	Callback: func(inputs []interface{}) (interface{}, error) {
		var list []string

		for _, input := range inputs {
			list = append(list, input.(string))
		}

		return strings.Join(list, ""), nil
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
			return funk.Contains(variablesMap, split[1]), nil
		} else if split[0] == "arg" {
			return funk.Contains(argsMap, split[1]), nil
		}

		return false, nil
	},
}

// getFunctions provides all the functions available in Hil language
func getFunctions() map[string]ast.Function {
	return map[string]ast.Function{
		"env":       env,
		"lower":     lowerCase,
		"upper":     upperCase,
		"snakeCase": snakeCase,
		"camelCase": camelCase,
		"reverse":   reverse,
		"shuffle":   shuffle,
		"wordCount": wordCount,
		"length":    length,
		"toString":  toString,
		"append":    appendFunc,
		"insert":    insert,
		"defined":   defined,
	}
}
