package logger

import (
	"fmt"
	"gopkg.ilharper.com/x/colette"
	"gopkg.ilharper.com/x/supcolor"
	"os"
	"regexp"
	"strconv"
)

var (
	// From https://github.com/chalk/ansi-regex
	regexpAnsi          = regexp.MustCompile(`[\033\0233][[\]()#;?]*(?:(?:(?:;[-a-zA-Z\d/#&.:=?%@~_]+)*|[a-zA-Z\d]+(?:;[-a-zA-Z\d/#&.:=?%@~_]*)*)?\007|(?:\d{1,4}(?:;\d{0,4})*)?[\dA-PR-TZcf-nq-uy=><~])`)
	regexpAnsiColor     = regexp.MustCompile(`\033\[[\d;]+m`)
	regexpAnsiColorTrue = regexp.MustCompile(`38;2;(\d+);(\d+);(\d+)`)
	regexpAnsiColor256  = regexp.MustCompile(`38;5;(\d+)`)
)

type colorAdapter struct {
	level int8
}

func newColorAdapter(target *os.File) *colorAdapter {
	return &colorAdapter{
		level: supcolor.SupColor(target),
	}
}

// Adapt 16 (basic) color, xterm(1) 256 color or true color
// to which [os.File] supports.
//
// # Format
//
// A color control sequence consists of 3 parts:
//
//     (CSI)[     (colors)     m
//       |            |        |
//       |            |        CSI Control End
//       |      Color Params
//     CSI Control Start
//
// Where color params has three modes:
//
//     16 (basic) color, level = 1:
//     (CSI)[     (foreground)     m
//                     |
//            Foreground: 30-37
//
//     xterm(1) 256 color, level = 2:
//     (CSI)[     38;5;(foreground)     m
//                          |
//                  Foreground: 0-255
//
//     true color, level = 3:
//     (CSI)[     38;2;R;G;B     m
//                       |
//           Foreground: (R, G, B)
//
// Mission of this function is to
// downgrade these control sequences to the level
// [os.File] supported.
//
// See [ANSI Escape Code].
//
// [ANSI Escape Code]: https://en.wikipedia.org/wiki/ANSI_escape_code
func (adapter *colorAdapter) adaptColor(s string) string {
	if adapter.level >= 3 {
		// This is an almighty terminal, just leave it to him
		return s
	}

	if adapter.level == 0 {
		// Terminal don't recognize ANSI codes... or even not a terminal
		// Just remove all sequences
		return regexpAnsi.ReplaceAllString(s, "")
	}

	// Terminal recognizes ANSI codes
	// Let's assume terminal can recognizes all codes except color
	// So just replace color codes to terminal's level
	// And first convert true colors to 256
	s = regexpAnsiColor.ReplaceAllStringFunc(s, func(code string) string {
		var err error
		matches := regexpAnsiColorTrue.FindStringSubmatch(code)
		if len(matches) < 4 {
			// Not true color code, just return
			return code
		}
		r, err := strconv.Atoi(matches[1])
		if err != nil {
			// ?? Not a number? Just return
			return code
		}
		g, err := strconv.Atoi(matches[2])
		if err != nil {
			return code
		}
		b, err := strconv.Atoi(matches[3])
		if err != nil {
			return code
		}
		c256 := colette.RgbTo256(byte(r), byte(g), byte(b))
		return fmt.Sprintf("\033[38;5;%dm", c256)
	})

	if adapter.level == 2 {
		// Terminal recognizes xterm(1) 256 colors
		// So it's now safe to return
		return s
	}

	// adapter.level == 1
	// Terminal recognizes 16 (basic) colors
	// So convert remaining colors to 16
	return regexpAnsiColor.ReplaceAllStringFunc(s, func(code string) string {
		matches := regexpAnsiColor256.FindStringSubmatch(code)
		if len(matches) < 2 {
			// Not 256 color code, just return
			return code
		}
		c256, err := strconv.Atoi(matches[1])
		if err != nil {
			// ?? Not a number? Just return
			return code
		}
		c16 := colette.Color256To16(c256)
		if c16 < 8 {
			c16 += 30
		} else {
			c16 += 82
		}
		return fmt.Sprintf("\033[%dm", c16)
	})
}
