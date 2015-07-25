package readline

import "fmt"

type Command func(rl *ReadLine, key Key)

var commands = map[string]Command{
	"beginning-of-line": func(rl *ReadLine, key Key) {
		rl.Pos = 0
	},
	"end-of-line": func(rl *ReadLine, key Key) {
		rl.Pos = len(rl.Text)
	},
	"forward-char": func(rl *ReadLine, key Key) {
		if rl.Pos < len(rl.Text) {
			rl.Pos++
		}
	},
	"backward-char": func(rl *ReadLine, key Key) {
		if rl.Pos > 0 {
			rl.Pos--
		}
	},

	"self-insert": func(rl *ReadLine, key Key) {
		rl.Insert(byte(key.Ch))
	},

	"kill-line": func(rl *ReadLine, key Key) {
		rl.Text = rl.Text[:rl.Pos]
	},
	"unix-line-discard": func(rl *ReadLine, key Key) {
		copy(rl.Text, rl.Text[rl.Pos:])
		rl.Text = rl.Text[:len(rl.Text)-rl.Pos]
		rl.Pos = 0
	},
}

func DefaultBindings() map[string]string {
	b := map[string]string{
		"C-a": "beginning-of-line",
		"C-e": "end-of-line",
		"C-f": "forward-char",
		"C-b": "backward-char",
		"M-f": "forward-word",
		"M-b": "backward-word",

		"C-k": "kill-line",
		"C-u": "unix-line-discard",
	}
	for ch := ' '; ch <= '~'; ch++ {
		b[fmt.Sprintf("%c", ch)] = "self-insert"
	}
	return b
}
