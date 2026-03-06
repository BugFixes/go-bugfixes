package middleware

import (
	"io"

	"github.com/bugfixes/go-bugfixes/internal/term"
)

var (
	// Normal colors
	nRed    = term.NRed
	nGreen  = term.NGreen
	nYellow = term.NYellow
	nCyan   = term.NCyan
	// Bright colors
	bRed     = term.BRed
	bGreen   = term.BGreen
	bYellow  = term.BYellow
	bBlue    = term.BBlue
	bMagenta = term.BMagenta
	bCyan    = term.BCyan
	bWhite   = term.BWhite

	reset = term.Reset
)

// IsTTY reports whether stdout appears to be a terminal.
var IsTTY bool

func init() {
	IsTTY = term.IsTTY
}

// cW writes a color-formatted string to w.
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	term.CW(w, useColor, color, s, args...)
}
