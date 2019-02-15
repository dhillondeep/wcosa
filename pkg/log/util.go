package log

// Generic Writeln function
func Writeln(args ...interface{}) bool {
	return Write(append(args, true)...)
}

// Shorthands
func Info(args ...interface{}) {
	Write(append(args, INFO)...)
}

func Infoln(args ...interface{}) {
	Writeln(append(args, INFO)...)
}

func Verb(args ...interface{}) {
	Write(append(args, VERB)...)
}

func Verbln(args ...interface{}) {
	Writeln(append(args, VERB)...)
}

func Warn(args ...interface{}) {
	Write(append(args, WARN, Yellow))
}

func Warnln(args ...interface{}) {
	Writeln(append(args, WARN, Yellow)...)
}

func Err(args ...interface{}) {
	Write(append(args, ERR, Red)...)
}

func Errln(args ...interface{}) {
	Writeln(append(args, ERR, Red)...)
}

func WriteSuccess(args ...interface{}) {
	Writeln(append(args, Green, "success")...)
}

func WriteFailure(args ...interface{}) {
	Writeln(append(args, Red, "failure")...)
}
