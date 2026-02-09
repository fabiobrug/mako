package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// mako-menu - Standalone menu for Mako shell
// Uses only /dev/tty for display, writes choice to stdout

type MenuItem struct {
	Label string
	Value string
}

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	title := os.Args[1]
	var items []MenuItem

	for i := 2; i < len(os.Args); i++ {
		parts := strings.SplitN(os.Args[i], "|", 2)
		if len(parts) == 2 {
			items = append(items, MenuItem{Label: parts[0], Value: parts[1]})
		}
	}

	if len(items) == 0 {
		os.Exit(1)
	}

	choice := showMenu(title, items)

	// Write choice to stdout for parent to capture
	fmt.Print(choice)
}

func showMenu(title string, items []MenuItem) string {
	// Open /dev/tty for ALL display I/O
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return items[0].Value
	}
	defer tty.Close()

	fd := tty.Fd()
	oldState, _ := getTermios(fd)
	makeRaw(fd)
	defer restoreTermios(fd, oldState)

	selected := 0

	// Shark colors
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	gray := "\033[38;2;150;150;150m"
	white := "\033[38;2;255;255;255m"
	reset := "\033[0m"

	menuLines := len(items) + 5
	
	draw := func() {
		// Draw menu
		tty.WriteString("\r\033[K\n") // Clear and add blank line
		tty.WriteString("  " + lightBlue + "╭─ " + white + title + reset + "\033[K\r\n")
		tty.WriteString("  " + lightBlue + "│" + reset + "\033[K\r\n")

		for i, item := range items {
			if i == selected {
				tty.WriteString("  " + lightBlue + "│" + reset + "  " + cyan + "❯ " + item.Label + reset + "\033[K\r\n")
			} else {
				tty.WriteString("  " + lightBlue + "│" + reset + "    " + gray + item.Label + reset + "\033[K\r\n")
			}
		}

		tty.WriteString("  " + lightBlue + "│" + reset + "\033[K\r\n")
		tty.WriteString("  " + lightBlue + "╰─" + reset + " " + gray + "Use ↑↓ arrows, Enter to select" + reset + "\033[K")
		// Don't add newline at end - keeps cursor on last line
	}
	
	redraw := func() {
		// Move up to start of menu (cursor is at end of last line)
		for i := 0; i < menuLines-1; i++ {
			tty.WriteString("\033[A")
		}
		// Now we're at start of first line, clear everything
		tty.WriteString("\033[J") // Clear from cursor to end of screen
		// Draw fresh menu
		draw()
	}

	// Hide cursor
	tty.WriteString("\033[?25l")
	defer tty.WriteString("\033[?25h")

	draw()

	// Input loop
	buf := make([]byte, 3)
	for {
		n, _ := tty.Read(buf[:1])
		if n == 0 {
			continue
		}

		switch buf[0] {
		case 13, 10: // Enter
			// Clear the menu properly
			// Move up to first line of menu (cursor is at end of last line)
			for i := 0; i < menuLines-1; i++ {
				tty.WriteString("\033[A")
			}
			// Move to start of line and clear everything below
			tty.WriteString("\r\033[J")
			return items[selected].Value

		case 3: // Ctrl+C
			// Clear the menu properly
			// Move up to first line of menu (cursor is at end of last line)
			for i := 0; i < menuLines-1; i++ {
				tty.WriteString("\033[A")
			}
			// Move to start of line and clear everything below
			tty.WriteString("\r\033[J")
			return "cancel"

		case 27: // ESC
			// Small delay to ensure the full escape sequence arrives
			time.Sleep(10 * time.Millisecond)
			// Read the next bytes for arrow keys
			n, _ = tty.Read(buf[1:3])
			if n >= 2 && buf[1] == 91 {
				switch buf[2] {
				case 65: // Up
					selected--
					if selected < 0 {
						selected = len(items) - 1
					}
					redraw()
				case 66: // Down
					selected++
					if selected >= len(items) {
						selected = 0
					}
					redraw()
				}
			}
		}
	}
}

func getTermios(fd uintptr) (*syscall.Termios, error) {
	t := &syscall.Termios{}
	_, _, e := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS,
		uintptr(unsafe.Pointer(t)), 0, 0, 0)
	if e != 0 {
		return nil, e
	}
	return t, nil
}

func makeRaw(fd uintptr) {
	t, _ := getTermios(fd)
	t.Lflag &^= syscall.ECHO | syscall.ICANON
	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0
	syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS,
		uintptr(unsafe.Pointer(t)), 0, 0, 0)
}

func restoreTermios(fd uintptr, t *syscall.Termios) {
	syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS,
		uintptr(unsafe.Pointer(t)), 0, 0, 0)
}
