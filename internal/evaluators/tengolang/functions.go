package tengolang

import (
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/stdlib"
	"github.com/huandu/xstrings"
)

// GetModelMap provides modules to be used inside tengo language
func GetModuleMap(names ...string) *objects.ModuleMap {
	modules := stdlib.GetModuleMap(names...)

	for _, name := range names {
		if mod := moduleMap[name]; mod != nil {
			modules.AddBuiltinModule(name, mod)
		}
	}

	return modules
}

var moduleMap = map[string]map[string]objects.Object{
	"wstrings": wUtilsModule,
}

var wUtilsModule = map[string]objects.Object{
	"snakeCase": &objects.UserFunction{Name: "snakeCase", Value: snakeCase},
	"camelCase": &objects.UserFunction{Name: "camelCase", Value: camelCase},
	"reverse":   &objects.UserFunction{Name: "reverse", Value: reverse},
	"shuffle":   &objects.UserFunction{Name: "shuffle", Value: shuffle},
	"wordCount": &objects.UserFunction{Name: "wordCount", Value: wordCount},
	"length":    &objects.UserFunction{Name: "length", Value: length},
}

func snakeCase(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	input, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return &objects.String{Value: xstrings.ToSnakeCase(input)}, nil
}

func camelCase(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	input, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return &objects.String{Value: xstrings.ToCamelCase(input)}, nil
}

func reverse(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	input, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return &objects.String{Value: xstrings.Reverse(input)}, nil
}

func shuffle(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	input, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return &objects.String{Value: xstrings.Shuffle(input)}, nil
}

func wordCount(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	input, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return &objects.Int{Value: int64(xstrings.WordCount(input))}, nil
}

func length(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}

	input, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return &objects.Int{Value: int64(xstrings.Len(input))}, nil
}
