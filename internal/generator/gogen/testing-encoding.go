package gogen

import (
	"fmt"
	"sort"
	"text/template"

	"github.com/sirkon/ch-encode/internal/generator"
)

// Int8TestEncoding ...
func (gg *GoGen) Int8TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "byte(%s),", source)
	return err
}

// Int16TestEncoding ...
func (gg *GoGen) Int16TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "int16(%s),", source)
	return err
}

// Int32TestEncoding ...
func (gg *GoGen) Int32TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "int32(%s),", source)
	return err
}

// Int64TestEncoding ...
func (gg *GoGen) Int64TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "int64(%s),", source)
	return err
}

// Uint8TestEncoding ...
func (gg *GoGen) Uint8TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "byte(%s),", source)
	return err
}

// Uint16TestEncoding ...
func (gg *GoGen) Uint16TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "uint16(%s),", source)
	return err
}

// Uint32TestEncoding ...
func (gg *GoGen) Uint32TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "uint32(%s),", source)
	return err
}

// Uint64TestEncoding ...
func (gg *GoGen) Uint64TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "uint64(%s),", source)
	return err
}

// Dec32TestEncoding ...
func (gg *GoGen) Dec32TestEncoding(scale int, source string) error {
	gg.regImport("", "github.com/sirkon/decconv")
	text := "decconv.Encode32(%d, int32(%s)),"
	_, err := fmt.Fprintf(gg.dest, text, scale, source)
	return err
}

// Dec64TestEncoding ...
func (gg *GoGen) Dec64TestEncoding(scale int, source string) error {
	gg.regImport("", "github.com/sirkon/decconv")
	text := "decconv.Encode64(%d, int64(%s)),"
	_, err := fmt.Fprintf(gg.dest, text, scale, source)
	return err
}

// Dec128TestEncoding ...
func (gg *GoGen) Dec128TestEncoding(scale int, source string) error {
	gg.regImport("", "github.com/sirkon/decconv")
	text := "decconv.Encode128(%d, %s.Lo, %s.Hi),"
	_, err := fmt.Fprintf(gg.dest, text, scale, source, source)
	return err
}

// Float32TestEncoding ...
func (gg *GoGen) Float32TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "float32(%s),", source)
	return err
}

// Float64TestEncoding ...
func (gg *GoGen) Float64TestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "float64(%s),", source)
	return err
}

// EnumTestEncoding ...
func (gg *GoGen) EnumTestEncoding(source string, safeMapping map[string]int) error {
	enumItems := NewEnumItemSlice(safeMapping)
	sort.Sort(enumItems)
	text :=
		`
        func() string {
            revMapping := map[int]string{
        {{ range .mapping }}    {{ .Value }}: "{{ .Key }}",
        {{ end }}
            }
            key, ok := revMapping[int({{ .source }})]
            if !ok {
                panic(fmt.Errorf("Value %d has no key mapped to", {{ .source }}));
            }
            return key
        }(),
        `
	tmpl, err := template.New("enum_test_encoding").Parse(text)
	if err != nil {
		return err
	}
	gg.regImport("", "fmt")
	err = tmpl.Execute(gg.dest, map[string]interface{}{
		"source":  source,
		"mapping": enumItems,
	})
	return err
}

// DateTestEncoding ...
func (gg *GoGen) DateTestEncoding(source string) error {
	text :=
		`
        func() string {
           timestamp := int64(%s)*3600*24
           tm := stdtime.Unix(timestamp, 0)
           moscowZone, _ := stdtime.LoadLocation("Europe/Moscow")
           moscow := tm.In(moscowZone)
           return moscow.Format("2006-01-02")
        }(),
        `
	_, err := fmt.Fprintf(gg.dest, text, source)
	return err
}

// DateTimeTestEncoding ...
func (gg *GoGen) DateTimeTestEncoding(source string) error {
	text :=
		`
        func() string {
           timestamp := int64(%s)
           tm := stdtime.Unix(timestamp, 0)
           moscowZone, _ := stdtime.LoadLocation("Europe/Moscow")
           moscow := tm.In(moscowZone)
           return moscow.Format("2006-01-02T15:04:05")
        }(),
        `
	_, err := fmt.Fprintf(gg.dest, text, source)
	return err

}

// StringTestEncoding ...
func (gg *GoGen) StringTestEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "string(%s),", source)
	return err
}

// FixedStringTestEncoding ...
func (gg *GoGen) FixedStringTestEncoding(source string, length int) error {
	_, err := fmt.Fprintf(gg.dest, "string(%s),", source)
	return err
}

// UUIDTestEncoding ...
func (gg *GoGen) UUIDTestEncoding(source string) error {
	gg.useGoogleUUID()
	_, err := fmt.Fprintf(gg.dest, "googleUUID.UUID(%s).String(),", source)
	return err
}

// ArrayTestEncoding ...
func (gg *GoGen) ArrayTestEncoding(source string, field generator.Field) error {
	text :=
		`
        func() (res []%s) {
           for _, arrayItem := range %s {
               res = append(res, `
	_, err := fmt.Fprintf(gg.dest, text, field.TestingTypeName(gg), source)
	if err != nil {
		return err
	}
	if err = field.TestEncoding("arrayItem", gg); err != nil {
		return err
	}
	if err = gg.RawData(");\n};\nreturn res;\n}(),\n"); err != nil {
		return err
	}
	return nil
}

// NullableTestEncoding ...
func (gg *GoGen) NullableTestEncoding(source string, field generator.Field) error {
	_, err := fmt.Fprintf(gg.dest, "(*%s)(%s),\n", field.TestingTypeName(gg), source)
	return err
}

// NullableStringTestEncoding
func (gg *GoGen) NullableStringTestEncoding(source string) error {
	text :=
		`
	   func() *string {
	       if %s == nil { return nil}
	       res := string(%s)
	       return &res
	   }(),`
	_, err := fmt.Fprintf(gg.dest, text, source, source)
	return err
}

// NullableArrayTestEncoding
func (gg *GoGen) NullableArrayTestEncoding(source string, field generator.Field) error {
	return gg.NullableTestEncoding(source, field)
}
