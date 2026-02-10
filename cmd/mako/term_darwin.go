//go:build darwin

package main

const (
	tcgets = 0x40487413 // TIOCGETA on macOS
	tcsets = 0x80487414 // TIOCSETA on macOS
)

func getTermiosCmd() uintptr {
	return tcgets
}

func setTermiosCmd() uintptr {
	return tcsets
}
