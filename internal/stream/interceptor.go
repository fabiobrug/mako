package stream

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fabiobrug/mako.git/internal/buffer"
	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/shell"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

type Interceptor struct {
	buffer *buffer.RingBuffer
	db     *database.DB
	writer io.Writer
}

func NewInterceptor(bufferSize int) *Interceptor {
	return &Interceptor{
		buffer: buffer.NewRingBuffer(bufferSize),
	}
}

func (i *Interceptor) SetDatabase(db *database.DB) {
	i.db = db
}

func (i *Interceptor) Tee(dst io.Writer, src io.Reader) error {
	i.writer = dst
	buf := make([]byte, 1024)
	lineBuffer := bytes.NewBuffer(nil)

	cmdFile := filepath.Join(os.Getenv("HOME"), ".mako", "last_command.txt")

	for {
		n, err := src.Read(buf)
		if n > 0 {
			data := buf[:n]

			// Check if this chunk contains the marker BEFORE writing to dst
			dataStr := string(data)
			if strings.Contains(dataStr, "<<<MAKO_EXECUTE>>>") {
				// Don't write the marker to output
				// Split on marker and only write the part before it
				parts := strings.Split(dataStr, "<<<MAKO_EXECUTE>>>")
				if len(parts) > 0 && parts[0] != "" {
					dst.Write([]byte(parts[0]))
				}

				// Clear the current line
				dst.Write([]byte("\r\033[K"))

				// Read and execute command
				if cmdBytes, errRead := os.ReadFile(cmdFile); errRead == nil {
					actualCommand := strings.TrimSpace(string(cmdBytes))

					shouldIntercept, output, cmdErr := shell.InterceptCommand(actualCommand, i.db)
					if shouldIntercept {
						if cmdErr != nil {
							errorMsg := fmt.Sprintf("Error: %v\r\n", cmdErr)
							dst.Write([]byte(errorMsg))
						} else if output != "" {
							// Ensure output has proper line endings
							output = strings.ReplaceAll(output, "\n", "\r\n")
							if !strings.HasSuffix(output, "\r\n") {
								output += "\r\n"
							}
							dst.Write([]byte(output))
						}
					}
					os.Remove(cmdFile)
				}

				// Write any remaining data after the marker
				if len(parts) > 1 {
					remaining := strings.Join(parts[1:], "")
					if remaining != "" {
						dst.Write([]byte(remaining))
					}
				}

				// Skip the normal line buffering for this chunk
				continue
			}

			// Normal path: write to destination
			dst.Write(data)

			// Buffer for line detection
			lineBuffer.Write(data)
			fullData := lineBuffer.Bytes()
			lastNewline := bytes.LastIndexByte(fullData, '\n')

			if lastNewline >= 0 {
				lines := bytes.Split(fullData[:lastNewline+1], []byte{'\n'})

				for _, line := range lines {
					lineStr := string(line)
					cleanLine := i.stripANSI(lineStr)
					cleanLine = strings.TrimSpace(cleanLine)

					if len(cleanLine) > 0 {
						i.buffer.Write(cleanLine)
					}
				}

				lineBuffer.Reset()
				if lastNewline+1 < len(fullData) {
					lineBuffer.Write(fullData[lastNewline+1:])
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				remaining := lineBuffer.Bytes()
				if len(remaining) > 0 {
					cleanLine := i.stripANSI(string(remaining))
					cleanLine = strings.TrimSpace(cleanLine)
					if len(cleanLine) > 0 {
						i.buffer.Write(cleanLine)
					}
				}
				return nil
			}
			return err
		}
	}
}

func (i *Interceptor) stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func (i *Interceptor) GetRecentLines(n int) []string {
	return i.buffer.GetLines(n)
}

func (i *Interceptor) GetAllLines() []string {
	return i.buffer.GetAll()
}

func (i *Interceptor) Clear() {
	i.buffer.Clear()
}
