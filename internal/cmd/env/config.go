package env

// store info about variables (readonly or not)
var envMeta = map[string]bool{
	"WIOROOT": true,
	"WIOOS":   true,
	"WIOARCH": true,
	"WIOPATH": true,
}
