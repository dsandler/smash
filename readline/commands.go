package readline

import (
	"fmt"
	"smash/keys"
)

type Command func(rl *ReadLine, key keys.Key)

var commands = map[string]Command{
	// Moving
	"beginning-of-line": func(rl *ReadLine, key keys.Key) {
		rl.Pos = 0
	},
	"end-of-line": func(rl *ReadLine, key keys.Key) {
		rl.Pos = len(rl.Text)
	},
	"forward-char": func(rl *ReadLine, key keys.Key) {
		if rl.Pos < len(rl.Text) {
			rl.Pos++
		}
	},
	"backward-char": func(rl *ReadLine, key keys.Key) {
		if rl.Pos > 0 {
			rl.Pos--
		}
	},

	// History
	"accept-line": func(rl *ReadLine, key keys.Key) {
		if rl.Accept != nil {
			rl.Pos = -1
			rl.Accept(rl.String())
		}
	},

	// Text
	"backward-delete-char": func(rl *ReadLine, key keys.Key) {
		if rl.Pos == 0 {
			return
		}
		copy(rl.Text[rl.Pos-1:], rl.Text[rl.Pos:])
		rl.Text = rl.Text[:len(rl.Text)-1]
		rl.Pos--
	},
	"self-insert": func(rl *ReadLine, key keys.Key) {
		rl.Insert(byte(key.Sym))
	},

	// Killing
	"kill-line": func(rl *ReadLine, key keys.Key) {
		rl.Text = rl.Text[:rl.Pos]
	},
	"unix-line-discard": func(rl *ReadLine, key keys.Key) {
		copy(rl.Text, rl.Text[rl.Pos:])
		rl.Text = rl.Text[:len(rl.Text)-rl.Pos]
		rl.Pos = 0
	},

	// Completion
	"complete": func(rl *ReadLine, key keys.Key) {
		rl.startComplete()
	},
}

func DefaultBindings() map[string]string {
	b := map[string]string{
		// Moving
		"C-a": "beginning-of-line",
		"C-e": "end-of-line",
		"C-f": "forward-char",
		"C-b": "backward-char",
		"M-f": "forward-word",
		"M-b": "backward-word",

		"Right": "forward-char",
		"Left":  "backward-char",

		// History
		"Enter": "accept-line",

		// Text
		"C-h":       "backward-delete-char",
		"Backspace": "backward-delete-char",

		// Killing
		"C-k": "kill-line",
		"C-u": "unix-line-discard",

		// Completion
		"Tab": "complete",
	}
	for ch := ' '; ch <= '~'; ch++ {
		b[fmt.Sprintf("%c", ch)] = "self-insert"
	}
	return b
}
