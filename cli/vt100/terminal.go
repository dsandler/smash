package vt100

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf8"
)

type Bits uint16

func (b Bits) Get(ofs uint, count uint) uint16 {
	return (uint16(b) >> ofs) & (uint16(1<<count) - 1)
}
func (b *Bits) Set(ofs uint, count uint, val uint) {
	mask := (uint16(1<<count) - 1) << ofs
	*b = Bits((uint16(*b) & ^uint16(mask)) | uint16(val)<<ofs)
}

// xxxx xxIB AAAA CCCC
//  I = inverse
//  B = bright
//  A = background color
//  C = foreground color
type Attr Bits

func (a Attr) Color() int {
	return int(Bits(a).Get(0, 4))
}
func (a *Attr) SetColor(color int) {
	(*Bits)(a).Set(0, 4, uint(color))
}

func (a Attr) Bright() bool {
	return Bits(a).Get(8, 1) != 0
}
func (a *Attr) SetBright(bright bool) {
	flag := uint(0)
	if bright {
		flag = 1
	}
	(*Bits)(a).Set(8, 1, flag)
}

func (a Attr) Inverse() bool {
	return Bits(a).Get(9, 1) != 0
}
func (a *Attr) SetInverse(inverse bool) {
	flag := uint(0)
	if inverse {
		flag = 1
	}
	(*Bits)(a).Set(9, 1, flag)
}

func (a Attr) BackColor() int {
	return int(Bits(a).Get(4, 4))
}
func (a *Attr) SetBackColor(color int) {
	(*Bits)(a).Set(4, 4, uint(color))
}

func (a Attr) Validate() error {
	if c := a.Color(); c < 0 || c > 8 {
		return fmt.Errorf("%s: bad color", a)
	}
	if c := a.BackColor(); c < 0 || c > 8 {
		return fmt.Errorf("%s: bad back color", a)
	}
	if uint16(a)&0xFC00 != 0 {
		return fmt.Errorf("%s: extra bits", uint16(a))
	}
	return nil
}

func (a Attr) String() string {
	fields := []string{}
	if a.Inverse() {
		fields = append(fields, "inverse")
	}
	if a.Bright() {
		fields = append(fields, "bright")
	}
	if fg := a.Color(); fg != 0 {
		fields = append(fields, fmt.Sprintf("fg:%d", fg))
	}
	if bg := a.BackColor(); bg != 0 {
		fields = append(fields, fmt.Sprintf("bg:%d", bg))
	}
	return fmt.Sprintf("Attr{%s}", strings.Join(fields, ","))
}

func showChar(ch byte) string {
	if ch >= ' ' && ch <= '~' {
		return fmt.Sprintf("'%c'", ch)
	} else {
		return fmt.Sprintf("%#x", ch)
	}
}

type Cell struct {
	Ch   rune
	Attr Attr
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell{%q, %s}", c.Ch, c.Attr)
}

type FeatureLog map[string]int

func (f FeatureLog) Add(text string, args ...interface{}) {
	if _, known := f[text]; !known {
		log.Printf("TODO: "+text, args...)
	}
	f[text]++
}

type Terminal struct {
	Title      string
	Lines      [][]Cell
	Width      int
	Height     int
	Input      io.Writer
	HideCursor bool
	TODOs      FeatureLog

	// Index of first displayed line; greater than 0 when content has
	// scrolled off the top of the terminal.
	Top int
	// The 0-based position of the cursor.
	Row, Col int
	Attr     Attr
}

func NewTerminal() *Terminal {
	return &Terminal{
		Lines:  make([][]Cell, 1),
		Width:  80,
		Height: 24,
		TODOs:  FeatureLog{},
	}
}

func promoteEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

func (t *Terminal) fixPosition() {
	if t.Row >= t.Top+t.Height {
		t.Top++
	}
	for t.Row >= len(t.Lines) {
		t.Lines = append(t.Lines, make([]Cell, 0))
	}
	for t.Col > len(t.Lines[t.Row]) {
		t.Lines[t.Row] = append(t.Lines[t.Row], Cell{' ', 0})
	}
}

func (t *Terminal) Read(r io.ByteScanner) error {
	c, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch {
	case c == 0x7: // bell
		// ignore
	case c == 0x8: // backspace
		if t.Col > 0 {
			t.Col--
		}
	case c == 0x1b:
		return t.readEscape(r)
	case c == '\r':
		t.Col = 0
	case c == '\n':
		t.Col = 0
		t.Row++
		t.fixPosition()
	case c == '\t':
		t.Col += 8 - (t.Col % 8)
		t.fixPosition()
	case c >= ' ' && c <= '~':
		t.writeRune(rune(c), t.Attr)
	default:
		r.UnreadByte()
		return t.readUTF8(r)
	}
	return nil
}

func (t *Terminal) writeRune(r rune, attr Attr) {
	if t.Col == t.Width {
		t.Row++
		t.Col = 0
	}
	t.Col++
	t.fixPosition()
	t.Lines[t.Row][t.Col-1] = Cell{r, attr}
}

func (t *Terminal) readUTF8(r io.ByteScanner) error {
	c, err := r.ReadByte()
	if err != nil {
		return err
	}

	attr := t.Attr

	var uc rune
	n := 0
	switch {
	case c&0xE0 == 0xB0:
		uc = rune(c & 0x1F)
		n = 2
	case c&0xF0 == 0xE0:
		uc = rune(c & 0x0F)
		n = 3
	default:
		if c&0xF0 == 0xF0 {
			log.Printf("term: not yet implemented: utf8 start %#v", c)
		}
		attr.SetInverse(true)
		t.writeRune('@', attr)
		return nil
	}

	for i := 1; i < n; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return err
		}
		if c&0xC0 != 0x80 {
			log.Printf("term: not yet implemented: utf8 continuation %#v", c)
			attr.SetInverse(true)
			uc = '@'
			break
		}
		uc = uc<<6 | rune(c&0x3F)
	}
	t.writeRune(uc, attr)
	return nil
}

func (t *Terminal) readEscape(r io.ByteScanner) error {
	// http://invisible-island.net/xterm/ctlseqs/ctlseqs.html
	c, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch {
	case c == '(':
		c, err := r.ReadByte()
		if err != nil {
			return err
		}
		switch c {
		case 'B': // US ASCII
			// ignore
		default:
			t.TODOs.Add("g0 charset %s", showChar(c))
		}
	case c == '=':
		t.TODOs.Add("application keypad")
	case c == '>':
		t.TODOs.Add("normal keypad")
	case c == '[':
		return t.readCSI(r)
	case c == ']':
		// OSC Ps ; Pt BEL
		n, err := t.readInt(r)
		if err != nil {
			return err
		}
		_, err = t.expect(r, ';')
		if err != nil {
			return err
		}
		text, err := t.readTo(r, 0x7)
		if err != nil {
			return err
		}
		switch n {
		case 0, 1, 2:
			t.Title = string(text)
		default:
			log.Printf("term: bad OSC %d", n)
		}
	case c == 'M': // move up/insert line
		if t.Row == 0 {
			// Insert line above.
			t.Lines = append(t.Lines, nil)
			copy(t.Lines[1:], t.Lines)
			t.Lines[0] = make([]Cell, 0)
		} else {
			if t.Row == t.Top {
				t.Top--
				if len(t.Lines) > t.Top+t.Height {
					t.Lines = t.Lines[:t.Top+t.Height-1]
				}
			}
			t.Row--
		}
	default:
		log.Printf("term: unknown escape %s", showChar(c))
	}
	return nil
}

func readArgs(args []int, values ...*int) {
	for i, val := range values {
		if i < len(args) {
			*val = args[i]
		}
	}
}

// mapColor converts a CSI color (e.g. 0=black, 1=red) to the term
// representation (0=default, 1=black).
func mapColor(color int, arg int) int {
	switch {
	case color == 8:
		log.Printf("term: bad color %d", arg)
		return 0
	case color == 9:
		return 0
	default:
		return color + 1
	}
}

func (t *Terminal) readCSI(r io.ByteScanner) error {
	// CSI
	var args []int

	qflag := false
	gtflag := false
L:
	c, err := r.ReadByte()
	if err != nil {
		return err
	}

	switch {
	case c >= '0' && c <= '9':
		r.UnreadByte()
		n, err := t.readInt(r)
		if err != nil {
			return err
		}
		args = append(args, n)

		c, err = r.ReadByte()
		if err != nil {
			return err
		}
		if c == ';' {
			goto L
		}
	case c == '?':
		qflag = true
		goto L
	case c == '>':
		gtflag = true
		goto L
	}

	switch {
	case c == '@': // insert blanks
		n := 1
		readArgs(args, &n)
		for i := 0; i < n; i++ {
			t.Lines[t.Row] = append(t.Lines[t.Row], Cell{})
		}
		copy(t.Lines[t.Row][t.Col+n:], t.Lines[t.Row][t.Col:])
		for i := 0; i < n; i++ {
			t.Lines[t.Row][t.Col+i] = Cell{' ', 0}
		}
	case c == 'A': // cursor up
		dy := 1
		readArgs(args, &dy)
		t.Row -= dy
		if t.Row < 0 {
			log.Printf("term: cursor up off top of screen?")
			t.Row = 0
		}
		t.fixPosition()
	case c == 'C': // cursor forward
		dx := 1
		readArgs(args, &dx)
		t.Col += dx
		t.fixPosition()
	case c == 'D': // cursor back
		dx := 1
		readArgs(args, &dx)
		t.Col -= dx
		t.fixPosition()
	case c == 'G': // cursor character absolute
		x := 1
		readArgs(args, &x)
		t.Col = x - 1
		t.fixPosition()
	case c == 'H': // move to position
		row := 1
		col := 1
		readArgs(args, &row, &col)
		t.Row = t.Top + row - 1
		t.Col = col - 1
		t.fixPosition()
	case c == 'J': // erase in display
		arg := 0
		readArgs(args, &arg)
		switch arg {
		case 0: // erase to end
			t.Lines = t.Lines[:t.Row+1]
			t.Lines[t.Row] = t.Lines[t.Row][:t.Col]
		case 2: // erase all
			t.Lines = t.Lines[:0]
			t.Row = 0
			t.Col = 0
			t.fixPosition()
		default:
			log.Printf("term: unknown erase in display %v", args)
		}
	case c == 'K': // erase in line
		arg := 0
		readArgs(args, &arg)
		switch arg {
		case 0: // erase to right
			t.Lines[t.Row] = t.Lines[t.Row][:t.Col]
		case 1:
			for i := 0; i < t.Col; i++ {
				t.Lines[t.Row][i] = Cell{' ', 0}
			}
		case 2:
			t.TODOs.Add("erase all line")
		default:
			log.Printf("term: unknown erase in line %v", args)
		}
	case c == 'L': // insert lines
		n := 1
		readArgs(args, &n)
		for i := 0; i < n; i++ {
			t.Lines = append(t.Lines, nil)
		}
		copy(t.Lines[t.Row+n:], t.Lines[t.Row:])
		for i := 0; i < n; i++ {
			t.Lines[t.Row+i] = make([]Cell, 0)
		}
	case c == 'P': // erase in line
		arg := 1
		readArgs(args, &arg)
		l := t.Lines[t.Row]
		if t.Col+arg > len(l) {
			arg = len(l) - t.Col
		}
		copy(l[t.Col:], l[t.Col+arg:])
		t.Lines[t.Row] = l[:len(l)-arg]
	case c == 'X': // erase characters
		t.TODOs.Add("erase characters %v", args)
	case !gtflag && c == 'c': // send device attributes (primary)
		t.TODOs.Add("send device attributes (primary) %v", args)
	case gtflag && c == 'c': // send device attributes (secondary)
		arg := 0
		readArgs(args, &arg)
		switch arg {
		case 0: // terminal id
			// ID is
			//   0 -> VT100
			//   0 -> firmware version 0
			//   0 -> always-zero param
			_, err := t.Input.Write([]byte("\x1b[0;0;0c"))
			return err
		default:
			t.TODOs.Add("send device attributes (secondary) %v", args)
		}
	case c == 'd': // line position
		arg := 1
		readArgs(args, &arg)
		t.Row = arg - 1
		t.fixPosition()
	case !qflag && (c == 'h' || c == 'l'): // reset mode
		reset := c == 'l'
		arg := 0
		readArgs(args, &arg)
		switch arg {
		default:
			t.TODOs.Add("reset mode %d %v", arg, reset)
		}
	case qflag && (c == 'h' || c == 'l'): // DEC private mode set/reset
		set := c == 'h'
		arg := 0
		readArgs(args, &arg)
		switch arg {
		case 1:
			t.TODOs.Add("application cursor keys mode")
		case 7: // wraparound mode
			t.TODOs.Add("wraparound mode")
		case 12: // blinking cursor
			// Ignore; this appears in cnorm/cvvis as a way to adjust the
			// "very visible cursor" state.
		case 25: // show cursor
			t.HideCursor = !set
		case 1049: // alternate screen buffer
			t.TODOs.Add("alternate screen buffer %v", set)
		default:
			log.Printf("term: unknown dec private mode %v %v", args, set)
		}
	case c == 'm': // reset
		if len(args) == 0 {
			args = append(args, 0)
		}
		for _, arg := range args {
			switch {
			case arg == 0:
				t.Attr = 0
			case arg == 1:
				t.Attr.SetBright(true)
			case arg == 7:
				t.Attr.SetInverse(true)
			case arg == 27:
				t.Attr.SetInverse(false)
			case arg >= 30 && arg < 40:
				t.Attr.SetColor(mapColor(arg-30, arg))
			case arg >= 40 && arg < 50:
				t.Attr.SetBackColor(mapColor(arg-40, arg))
			default:
				log.Printf("term: unknown color %v", args)
			}
		}
	case gtflag && c == 'n': // disable modifiers
		arg := 2
		readArgs(args, &arg)
		switch arg {
		case 0:
			t.TODOs.Add("disable modify keyboard")
		case 1:
			t.TODOs.Add("disable modify cursor keys")
		case 2:
			t.TODOs.Add("disable modify function keys")
		case 4:
			t.TODOs.Add("disable modify other keys")
		}
	case c == 'n': // device status report
		arg := 0
		readArgs(args, &arg)
		switch arg {
		case 5:
			_, err := t.Input.Write([]byte("\x1b[0n"))
			return err
		case 6:
			pos := fmt.Sprintf("\x1b[%d;%dR", t.Row+1, t.Col+1)
			_, err := t.Input.Write([]byte(pos))
			return err
		default:
			log.Printf("term: unknown status report arg %v", args)
		}
	case c == 'r': // set scrolling region
		top, bot := 1, 1
		readArgs(args, &top, &bot)
		if top == 1 && bot == t.Height {
			// Just setting the current region as scroll.
		} else {
			t.TODOs.Add("set scrolling region %v", args)
		}
	default:
		log.Printf("term: unknown CSI %v %s", args, showChar(c))
	}
	return nil
}

func (t *Terminal) expect(r io.ByteScanner, exp byte) (bool, error) {
	c, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	ok := c == exp
	if !ok {
		log.Printf("expect %s failed, got %s", showChar(exp), showChar(c))
	}
	return ok, nil
}

func (t *Terminal) readInt(r io.ByteScanner) (int, error) {
	n := 0
	for i := 0; i < 20; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return -1, err
		}
		if c >= '0' && c <= '9' {
			n = n*10 + int(c) - '0'
		} else {
			r.UnreadByte()
			return n, err
		}
	}
	return -1, fmt.Errorf("term: readInt overlong")
}

func (t *Terminal) readTo(r io.ByteScanner, end byte) ([]byte, error) {
	var buf []byte
	for i := 0; i < 1000; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c == end {
			return buf, nil
		}
		buf = append(buf, c)
	}
	return nil, fmt.Errorf("term: readTo(%s) overlong", showChar(end))
}

func (t *Terminal) DisplayString(input string) {
	r := strings.NewReader(input)
	var err error
	for err == nil {
		err = t.Read(r)
	}
}

func (t *Terminal) ToString() string {
	var buf [6]byte
	str := ""
	for _, l := range t.Lines {
		if str != "" {
			str += "\n"
		}
		for _, c := range l {
			n := utf8.EncodeRune(buf[:], c.Ch)
			str += string(buf[:n])
		}
	}
	return str
}

func (t *Terminal) Validate() error {
	for row, l := range t.Lines {
		for col, c := range l {
			if err := c.Attr.Validate(); err != nil {
				return fmt.Errorf("%d:%d: %s", row, col, err)
			}
		}
	}
	return nil
}
