package shell

import (
	"os"
	"strings"
)

// readLineFromTTY reads a line of input from /dev/tty with the prefilled text
func readLineFromTTY(prefill string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", err
	}
	defer tty.Close()

	// Write prefilled text
	tty.WriteString(prefill)

	// Read until newline - initialize with prefilled text
	result := []byte(prefill)
	buf := make([]byte, 1)
	
	for {
		n, err := tty.Read(buf)
		if err != nil {
			return "", err
		}
		if n > 0 {
			if buf[0] == '\n' || buf[0] == '\r' {
				break
			}
			// Handle backspace
			if buf[0] == 127 || buf[0] == 8 {
				if len(result) > 0 {
					result = result[:len(result)-1]
					// Visual backspace: move back, print space, move back again
					tty.WriteString("\b \b")
				}
			} else {
				result = append(result, buf[0])
				// Echo the character
				tty.Write(buf)
			}
		}
	}
	
	return string(result), nil
}

// wrapLine wraps a line of text to fit within maxWidth characters
func wrapLine(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := ""
	for _, word := range words {
		// If adding this word would exceed maxWidth, start a new line
		if len(currentLine)+len(word)+1 > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Single word is longer than maxWidth, split it
				lines = append(lines, word[:maxWidth])
				currentLine = word[maxWidth:]
			}
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
