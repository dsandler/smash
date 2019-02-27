package vt100

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeInput(input string) *bufio.Reader {
	buf := strings.NewReader(input)
	return bufio.NewReader(buf)
}

func mustRun(t *testing.T, term *Terminal, input string) {
	r := makeInput(input)
	var err error
	for err == nil {
		err = term.Read(r)
	}
	assert.Equal(t, err, io.EOF)
}

func assertPos(t *testing.T, term *Terminal, row, col int) {
	assert.Equal(t, row, term.Row)
	assert.Equal(t, col, term.Col)
}

func TestBasic(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "test")
	assert.Equal(t, "test", term.ToString())
	mustRun(t, term, "\nbar")
	assert.Equal(t, "test\nbar", term.ToString())
	mustRun(t, term, "\rfoo")
	assert.Equal(t, "test\nfoo", term.ToString())
	mustRun(t, term, "\n\n")
	assert.Equal(t, "test\nfoo\n\n", term.ToString())
	mustRun(t, term, "x\ty")
	assert.Equal(t, "test\nfoo\n\nx       y", term.ToString())
}

func TestTitle(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x1b]0;title\x07text")
	assert.Equal(t, "title", term.Title)
	assert.Equal(t, "text", term.ToString())
}

func TestReset(t *testing.T) {
	term := NewTerminal()
	term.Attr = 43
	mustRun(t, term, "\x1b[0m")
	assert.Equal(t, Attr(0), term.Attr)
	assert.Equal(t, "", term.ToString())
}

func TestColor(t *testing.T) {
	term := NewTerminal()
	assert.Equal(t, false, term.Attr.Bright())
	assert.Equal(t, false, term.Attr.Inverse())
	assert.Equal(t, 0, term.Attr.Color())

	mustRun(t, term, "\x1b[1;34m")
	assert.Equal(t, true, term.Attr.Bright())
	assert.Equal(t, 5, term.Attr.Color())
	assert.Equal(t, "", term.ToString())

	mustRun(t, term, "\x1b[7m")
	assert.Equal(t, true, term.Attr.Inverse())

	mustRun(t, term, "\x1b[m")
	assert.Equal(t, Attr(0), term.Attr)
}

func TestBackspace(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x08")
	assert.Equal(t, "", term.ToString())
	mustRun(t, term, "x\x08")
	assert.Equal(t, "x", term.ToString())
	mustRun(t, term, "ab\x08c")
	assert.Equal(t, "ac", term.ToString())
}

func TestEraseLine(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "hello")
	term.Col -= 2
	mustRun(t, term, "\x1b[K")
	assert.Equal(t, "hel", term.ToString())
	mustRun(t, term, "\x1b[1K")
	assert.Equal(t, "   ", term.ToString())
}

func TestEraseDisplay(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "hellofoo\b\b\b")
	mustRun(t, term, "\x1b[J")
	assert.Equal(t, "hello", term.ToString())
	mustRun(t, term, "\x1b[2J")
	assert.Equal(t, "", term.ToString())
}

func TestDelete(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "abcdef\x08\x08\x08\x1b[1P")
	assert.Equal(t, "abcef", term.ToString())

	// Check deleting past the end of the line.
	mustRun(t, term, "\x1b[5P")
	assert.Equal(t, "abc", term.ToString())
}

func TestBell(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x07")
	// ignored
	assert.Equal(t, "", term.ToString())
}

func TestPrivateModes(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x1b[?1049h")
	// ignored
	assert.Equal(t, "", term.ToString())

	mustRun(t, term, "\x1b[?7h")
	// ignored
	assert.Equal(t, "", term.ToString())
}

func TestScrollingRegion(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x1b[1;24r")
	// ignored
	assert.Equal(t, "", term.ToString())
}

func TestResetMode(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x1b[4l")
	// ignored
	assert.Equal(t, "", term.ToString())
}

func TestMoveTo(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "hello\x1b[HX")
	assert.Equal(t, "Xello", term.ToString())
	mustRun(t, term, "\x1b[1;3HX")
	assert.Equal(t, "XeXlo", term.ToString())
}

func TestMoveToLine(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "hello\n\n\x1b[2dfoo")
	assert.Equal(t, "hello\nfoo\n", term.ToString())
}

func TestCursor(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "foo\nbar")
	assertPos(t, term, 1, 3)

	mustRun(t, term, "\x1b[C")
	assertPos(t, term, 1, 4)
	mustRun(t, term, "\x1b[2C")
	assertPos(t, term, 1, 6)
	assert.Equal(t, "foo\nbar   ", term.ToString())

	mustRun(t, term, "\x1b[A!")
	assertPos(t, term, 0, 7)
	assert.Equal(t, "foo   !\nbar   ", term.ToString())

	mustRun(t, term, "\x1b[5D")
	assertPos(t, term, 0, 2)
}

func TestScrollUp(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "aaa\nbbb\nccc\n")
	assertPos(t, term, 3, 0)
	mustRun(t, term, "\x1bMX")
	assert.Equal(t, "aaa\nbbb\nXcc\n", term.ToString())
	mustRun(t, term, "\x1bMY")
	mustRun(t, term, "\x1bMZ")
	assert.Equal(t, "aaZ\nbYb\nXcc\n", term.ToString())
	mustRun(t, term, "\x1bM1")
	assert.Equal(t, "   1\naaZ\nbYb\nXcc\n", term.ToString())
}

func TestScrollUpDropLines(t *testing.T) {
	term := NewTerminal()
	term.Height = 3
	mustRun(t, term, "aaa\nbbb\nccc\n")
	assert.Equal(t, "aaa\nbbb\nccc\n", term.ToString())
	mustRun(t, term, "\x1bM\x1bM\x1bM\x1bMx")
	assert.Equal(t, "x\naaa\nbbb", term.ToString())
}

func TestWrap(t *testing.T) {
	term := NewTerminal()
	term.Width = 5
	mustRun(t, term, "1234567890")
	assert.Equal(t, "12345\n67890", term.ToString())
}

func TestUTF8(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\xe2\x96\xbd")
	assert.Equal(t, rune(0x25bd), term.Lines[0][0].Ch)
}

func TestStatusReport(t *testing.T) {
	term := NewTerminal()
	buf := &bytes.Buffer{}
	term.Input = buf
	mustRun(t, term, "\x1b[5n")
	assert.Equal(t, "", term.ToString())
	assert.Equal(t, "\x1b[0n", buf.String())

	buf.Reset()
	mustRun(t, term, "\x1b[6n")
	assert.Equal(t, "\x1b[1;1R", buf.String())
}

func TestCSIDisableModifiers(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x1b[>0n")
	assert.Equal(t, "", term.ToString())
	// TODO: implement the disabling, whatever that is.
}

func TestSendDeviceAttributes(t *testing.T) {
	term := NewTerminal()
	buf := &bytes.Buffer{}
	term.Input = buf
	mustRun(t, term, "\x1b[c")
	assert.Equal(t, "", term.ToString())
	assert.Equal(t, "", buf.String())
	mustRun(t, term, "\x1b[>c")
	assert.Equal(t, "", term.ToString())
	assert.Equal(t, "\x1b[0;0;0c", buf.String())
}

func TestHideCursor(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "\x1b[?25l")
	assert.Equal(t, true, term.HideCursor)
	mustRun(t, term, "\x1b[?25h")
	assert.Equal(t, false, term.HideCursor)
}

func TestInsertBlanks(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "ABC\b\b\x1b[@x")
	assert.Equal(t, "AxBC", term.ToString())
	mustRun(t, term, "\x1b[2@y")
	assert.Equal(t, "Axy BC", term.ToString())
}

func TestInsertLine(t *testing.T) {
	term := NewTerminal()
	mustRun(t, term, "foo\nbar\nbaz\n")
	mustRun(t, term, "\x1b[2A\x1b[L") // two lines up, insert line
	mustRun(t, term, "\nX")
	assert.Equal(t, "foo\n\nXar\nbaz\n", term.ToString())
}

func TestBinary(t *testing.T) {
	term := NewTerminal()
	// Don't choke on non-UTF8 inputs.
	// TODO: maybe render them with some special character to represent
	// mojibake.
	mustRun(t, term, "\xc8\x00\x64\x00")
	assert.Equal(t, "@@d@", term.ToString())
}

func TestAllColors(t *testing.T) {
	buf := &bytes.Buffer{}
	for i := 30; i < 50; i++ {
		fmt.Fprintf(buf, "\x1b[%dmx", i)
	}
	term := NewTerminal()
	mustRun(t, term, buf.String())
	x20 := "xxxxxxxxxx" + "xxxxxxxxxx"
	assert.Equal(t, x20, term.ToString())
	assert.Nil(t, term.Validate())
}
