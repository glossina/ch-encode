package gogen

import (
	"fmt"
	"text/template"

	"io"

	"github.com/sirkon/ch-encode/internal/generator"
)

// Int8Encoding ...
func (gg *GoGen) Int8Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Byte(byte(%s)));\n", source)
	return err
}

// Int16Encoding ...
func (gg *GoGen) Int16Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Int16(int16(%s)));\n", source)
	return err
}

// Int32Encoding ...
func (gg *GoGen) Int32Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Int32(int32(%s)));\n", source)
	return err
}

// Int64Encoding ...
func (gg *GoGen) Int64Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Int64(int64(%s)));\n", source)
	return err
}

// Uint8Encoding ...
func (gg *GoGen) Uint8Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Byte(byte(%s)));\n", source)
	return err
}

// Uint16Encoding ...
func (gg *GoGen) Uint16Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint16(uint16(%s)));\n", source)
	return err
}

// Uint32Encoding ...
func (gg *GoGen) Uint32Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint32(uint32(%s)));\n", source)
	return err
}

// Uint64Encoding ...
func (gg *GoGen) Uint64Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint64(uint64(%s)));\n", source)
	return err
}

func (gg *GoGen) Dec128Encoding(source string) error {
	_, err := fmt.Fprintf(gg.dest,
		""+
			"enc.buffer.Write(enc.helper.Uint64(uint64(%s.Lo)));\n"+
			"enc.buffer.Write(enc.helper.Uint64(uint64(%s.Hi)));\n",
		source, source)
	return err
}

// Float32Encoding ...
func (gg *GoGen) Float32Encoding(source string) error {
	gg.useUnsafe()
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint32(*(*uint32)(unsafe.Pointer(&%s))));\n", source)
	return err
}

// Float64Encoding ...
func (gg *GoGen) Float64Encoding(source string) error {
	gg.useUnsafe()
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint64(*(*uint64)(unsafe.Pointer(&%s))));\n", source)
	return err
}

// DateEncoding ...
func (gg *GoGen) DateEncoding(source string) error {
	gg.useTime()
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint16(uint16(%s)));\n", source)
	return err
}

// DateTimeEncoding ...
func (gg *GoGen) DateTimeEncoding(source string) error {
	gg.useTime()
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(enc.helper.Uint32(uint32(%s)));\n", source)
	return err
}

// StringEncoding ...
func (gg *GoGen) StringEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest,
		""+
			"enc.buffer.Write(enc.helper.Uleb128(uint32(len([]byte(%s)))));\n"+
			"enc.buffer.Write([]byte(%s));\n",
		source, source)
	return err
}

// FixedStringEncoding ...
func (gg *GoGen) FixedStringEncoding(source string, length int) error {
	gg.useFmt()
	text := `
if len({{.var}}) != {{.length}} {
    return fmt.Errorf("string {{.var}} must be {{.length}} bytes long, got %d bytes intsead (\"\033[1m%s\033[0m\", %v)", len({{.var}}), string({{.var}}), {{.var}})
}
enc.buffer.Write([]byte({{.var}}))
`
	t, err := template.New("fixed_string_encoding").Parse(text)
	if err != nil {
		return err
	}
	err = t.Execute(gg.dest, map[string]interface{}{
		"var":    source,
		"length": length,
	})
	return err
}

// UUIDEncoding ...
func (gg *GoGen) UUIDEncoding(source string) error {
	_, err := fmt.Fprintf(gg.dest, "enc.buffer.Write(%s[:]);\n", source)
	return err
}

// ArrayEncoding ...
func (gg *GoGen) ArrayEncoding(source string, field generator.Field) error {
	text := `
            enc.buffer.Write(enc.helper.Uleb128(uint32(len({{.var}}))));
            for _, arrayItem := range {{.var}} {
            `
	tmpl, err := template.New("array_encoding").Parse(text)
	if err != nil {
		return err
	}
	err = tmpl.Execute(gg.dest, map[string]string{
		"var": source,
	})
	if err != nil {
		return err
	}
	if err = field.Encoding("arrayItem", gg); err != nil {
		return err
	}
	return gg.RawData("}")
}

// NullableEncoding ...
func (gg *GoGen) NullableEncoding(source string, field generator.Field) error {
	if _, err := fmt.Fprintf(gg.dest, ";if %s != nil {\n", source); err != nil {
		return err
	}
	if _, err := io.WriteString(gg.dest, ";enc.buffer.WriteByte(byte(0));"); err != nil {
		return err
	}
	if err := field.Encoding("*"+source, gg); err != nil {
		return err
	}
	return gg.RawData("} else { enc.buffer.WriteByte(byte(1)) }")
}

// NullableArrayEncoding ...
func (gg *GoGen) NullableArrayEncoding(source string, field generator.Field) error {
	return gg.NullableEncoding(source, field)
}

// NullableStringEncoding
func (gg *GoGen) NullableStringEncoding(source string) error {
	if _, err := fmt.Fprintf(gg.dest, ";if %s != nil {\n", source); err != nil {
		return err
	}
	_, err := fmt.Fprintf(gg.dest, ""+"enc.buffer.WriteByte(byte(0));\n")
	if err != nil {
		return err
	}
	if err := gg.StringEncoding(source); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(gg.dest, "; } else { enc.buffer.WriteByte(byte(1)) }"); err != nil {
		return err
	}
	return nil
}
