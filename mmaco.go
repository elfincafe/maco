package mmaco

import (
	"bytes"
	"reflect"
	"regexp"
	"time"
)

const (
	tagName     = "mmaco"
	helpCmdName = "help"
	trimSpace   = " \t\v\r\n\f"

	Unknown Kind = iota
	String
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	Time

	ShortOpt ArgType = iota
	LongOpt
	Value
	LongOptValue
)

type (
	Kind    int
	ArgType int

	Command struct {
		Name        string
		loc         *time.Location
		subCmds     map[string]*subCommand
		scOrder     []string
		opts        []*option
		start       time.Time
		startSubCmd time.Time
		endSubCmd   time.Time
		report      bool `mmaco:"short=r,long=report,desc=report verbosely."`
		help        bool `mmaco:"short=h,long=help,desc=this help."`
	}

	SubCommandInterface interface {
		Init()
		Validate() error
		Run([]string) error
	}

	subCommand struct {
		Name string
		Desc string
		cmd  reflect.Value
		opts []*option
		// hasValidate bool
		loc *time.Location
	}

	option struct {
		value     reflect.Value
		field     reflect.StructField
		loc       *time.Location
		specified bool
		Kind      Kind
		Name      string
		Short     string
		Long      string
		Required  bool
		Desc      string
		Default   string
		Format    string
		Handler   string
	}

	help struct {
		Desc string
	}
)

func toSnakeCase(s string) string {
	name := []byte{}
	for _, c := range []byte(s) {
		if c > 64 && c < 91 {
			name = append(name, byte(95), c+32)
		} else {
			name = append(name, c)
		}
	}
	name = regexp.MustCompile(`_{2,}`).ReplaceAll(name, []byte{95})
	return string(bytes.Trim(name, "_"))
}

func isAlphaNumeric(b byte) bool {
	if b >= 48 && b <= 57 { // 0 - 9
		return true
	} else if b >= 65 && b <= 90 { // A - Z
		return true
	} else if b >= 97 && b <= 122 { // a - z
		return true
	} else {
		return false
	}
}

func getFieldKind(field reflect.StructField) Kind {
	switch field.Type.Kind() {
	case reflect.Bool:
		return Bool
	case reflect.Int:
		return Int
	case reflect.Int8:
		return Int8
	case reflect.Int16:
		return Int16
	case reflect.Int32:
		return Int32
	case reflect.Int64:
		return Int64
	case reflect.Uint:
		return Uint
	case reflect.Uint8:
		return Uint8
	case reflect.Uint16:
		return Uint16
	case reflect.Uint32:
		return Uint32
	case reflect.Uint64:
		return Uint64
	case reflect.Float32:
		return Float32
	case reflect.Float64:
		return Float64
	case reflect.String:
		return String
	case reflect.Struct:
		if field.Type.PkgPath() == "time" && field.Type.Name() == "Time" {
			return Time
		} else {
			return Unknown
		}
	default:
		return Unknown
	}
}

func tryParseDateTime(s string) (time.Time, error) {
	layouts := []string{
		time.DateTime,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateOnly,
		time.TimeOnly,
		time.Kitchen,
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return t, err
}

// func isSubCommand(t reflect.Type) bool {
// 	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
// 		return false
// 	}
// 	if !hasRunMethod(t) {
// 		return false
// 	}
// 	return true
// }

// func hasInitMethod(t reflect.Type) bool {
// 	method, ok := t.MethodByName("Init")
// 	if !ok {
// 		return false
// 	}
// 	if method.Type.NumIn() != 1 {
// 		return false
// 	}
// 	if method.Type.NumOut() != 0 {
// 		return false
// 	}
// 	return true
// }

// func hasValidateMethod(t reflect.Type) bool {
// 	method, ok := t.MethodByName("Validate")
// 	if !ok {
// 		return false
// 	}
// 	if method.Type.NumIn() != 1 {
// 		return false
// 	}
// 	if method.Type.NumOut() != 1 || method.Type.Out(0).Kind() != reflect.Interface {
// 		return false
// 	}
// 	return true
// }

// func hasRunMethod(t reflect.Type) bool {
// 	method, ok := t.MethodByName("Run")
// 	if !ok {
// 		return false
// 	}
// 	if method.Type.NumIn() != 2 || method.Type.In(1).Kind() != reflect.Slice {
// 		return false
// 	}
// 	if method.Type.NumOut() != 1 || method.Type.Out(0).Kind() != reflect.Interface {
// 		return false
// 	}
// 	return true
// }

func hasDescField(t reflect.Type) bool {
	field, ok := t.FieldByName("Desc")
	if !ok {
		return false
	}
	if field.Type.Kind() != reflect.String {
		return false
	}
	return true
}
