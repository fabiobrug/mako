package ui

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

type MenuItem struct {
	Label string
	Value string
}

// ShowMenu displays an interactive menu with arrow key navigation
// Uses /dev/tty directly to avoid conflicts with PTY stream
func ShowMenu(title string, items []MenuItem) string {
	// Open /dev/tty for direct terminal access
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// Fallback to first option if tty fails
		return items[0].Value
	}
	defer tty.Close()

	// Save current terminal state
	oldState, err := getTermios(tty.Fd())
	if err != nil {
		return items[0].Value
	}

	// Set terminal to raw mode for menu
	if err := makeRaw(tty.Fd()); err != nil {
		return items[0].Value
	}
	defer restoreTermios(tty.Fd(), oldState)

	selected := 0

	// Color definitions
	cyan := "\033[38;2;93;173;226m"
	lightBlue := "\033[38;2;0;209;255m"
	gray := "\033[38;2;150;150;150m"
	reset := "\033[0m"

	// Helper to write to tty
	write := func(s string) {
		io.WriteString(tty, s)
	}

	// Initial render
	renderMenu := func() {
		write(fmt.Sprintf("\r\n%s%s%s\r\n", cyan, title, reset))
		for i, item := range items {
			if i == selected {
				write(fmt.Sprintf("  %s❯ %s%s\r\n", lightBlue, item.Label, reset))
			} else {
				write(fmt.Sprintf("  %s  %s%s\r\n", gray, item.Label, reset))
			}
		}
		write(fmt.Sprintf("\r\n%sUse ↑↓ arrows, Enter to select%s\r\n", gray, reset))
	}

	// Hide cursor
	write("\033[?25l")
	defer write("\033[?25h")

	renderMenu()

	// Read input byte by byte
	buf := make([]byte, 3) // Read up to 3 bytes for escape sequences
	for {
		n, err := tty.Read(buf[:1])
		if err != nil || n == 0 {
			break
		}

		switch buf[0] {
		case 13, 10: // Enter or newline
			// Clear menu
			clearLines(tty, len(items)+3)
			return items[selected].Value

		case 3: // Ctrl+C
			clearLines(tty, len(items)+3)
			return "cancel"

		case 27: // ESC - potential arrow key
			// Try to read the next 2 bytes (non-blocking style)
			n, err := tty.Read(buf[1:3])
			if err != nil || n < 2 {
				continue
			}

			// Check for arrow keys: ESC [ A/B/C/D
			if buf[1] == 91 { // '['
				oldSelected := selected
				switch buf[2] {
				case 65: // Up arrow
					selected = (selected - 1 + len(items)) % len(items)
				case 66: // Down arrow
					selected = (selected + 1) % len(items)
				}

				if oldSelected != selected {
					// Redraw menu
					clearLines(tty, len(items)+3)
					renderMenu()
				}
			}
		}
	}

	// Fallback
	clearLines(tty, len(items)+3)
	return items[selected].Value
}

func clearLines(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		io.WriteString(w, "\033[1A") // Move up
		io.WriteString(w, "\033[2K") // Clear line
		io.WriteString(w, "\r")      // Carriage return
	}
}

func getTermios(fd uintptr) (*syscall.Termios, error) {
	termios := &syscall.Termios{}
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCGETS,
		uintptr(unsafe.Pointer(termios)),
		0, 0, 0,
	)
	if errno != 0 {
		return nil, errno
	}
	return termios, nil
}

func makeRaw(fd uintptr) error {
	termios, err := getTermios(fd)
	if err != nil {
		return err
	}

	// Make raw
	termios.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	termios.Oflag &^= syscall.OPOST
	termios.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON |
		syscall.ISIG | syscall.IEXTEN
	termios.Cflag &^= syscall.CSIZE | syscall.PARENB
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	return restoreTermios(fd, termios)
}

func restoreTermios(fd uintptr, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCSETS,
		uintptr(unsafe.Pointer(termios)),
		0, 0, 0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}
