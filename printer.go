package log

import (
	"fmt"
	"io"
	"strings"
)

// BufferWriter is implemented by several very common writers such as
// bytes.Buffer, strings.Builder, and bufio.Writer.
type BufferWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
	fmt.Stringer
}

// Printer outputs string data using standard methods and panics on failure. It's
// typically intended for adding data to buffers where the only possible error is
// out of memory, which is a panic situation in any case.
type Printer struct {
	BufferWriter
}

// NewPrinter creates a new Printer instance on a default strings.Builder{}.
func NewPrinter() Printer {
	return Printer{&strings.Builder{}}
}

// NewPrinterOn creates a new Printer instance on the given BufferWriter (typically
// a bytes.Buffer, strings.Builder, or a bufio.Writer).
func NewPrinterOn(w BufferWriter) Printer {
	return Printer{w}
}

// Printf formats according to a format specifier and writes to p. It returns the number of bytes written.
func (p Printer) Printf(f string, as ...any) int {
	n, err := fmt.Fprintf(p, f, as...)
	if err != nil {
		panic(err)
	}
	return n
}

// Println formats using the default formats for its operands and writes to p. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes written.
func (p Printer) Println(f string) int {
	n, err := fmt.Fprintln(p, f)
	if err != nil {
		panic(err)
	}
	return n
}

// PrintString writes the contents of the string s to p.
func (p Printer) PrintString(s string) int {
	n, err := p.WriteString(s)
	if err != nil {
		panic(err)
	}
	return n
}

// PrintByte writes the byte b to p.
func (p Printer) PrintByte(b byte) {
	err := p.WriteByte(b)
	if err != nil {
		panic(err)
	}
}
