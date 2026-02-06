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

			// Always pass through to screen first
			if _, writeErr := dst.Write(data); writeErr != nil {
				return writeErr
			}

			// Process for buffer and command detection
			lineBuffer.Write(data)
			fullData := lineBuffer.Bytes()
			lastNewline := bytes.LastIndexByte(fullData, '\n')

			if lastNewline >= 0 {
				lines := bytes.Split(fullData[:lastNewline+1], []byte{'\n'})
				for _, line := range lines {
					if len(line) == 0 {
						continue
					}

					lineStr := string(line)
					cleanLine := i.stripANSI(lineStr)
					cleanLine = strings.TrimSpace(cleanLine)

					if len(cleanLine) > 0 {
						i.buffer.Write(cleanLine)

						// Check for execution marker
						if strings.Contains(cleanLine, "<<<MAKO_EXECUTE>>>") {
							// Read the command from file
							if cmdBytes, err := os.ReadFile(cmdFile); err == nil {
								actualCommand := strings.TrimSpace(string(cmdBytes))

								shouldIntercept, output, cmdErr := shell.InterceptCommand(actualCommand, i.db)
								if shouldIntercept {
									// Clear the marker line (move up one line, then clear it)
									dst.Write([]byte("\033[1A\r\033[K"))

									if cmdErr != nil {
										errorMsg := fmt.Sprintf("Error: %v\n", cmdErr)
										errorMsg = strings.ReplaceAll(errorMsg, "\n", "\r\n")
										dst.Write([]byte(errorMsg))
									} else if output != "" {
										// Replace \n with \r\n for proper terminal output
										output = strings.ReplaceAll(output, "\n", "\r\n")
										dst.Write([]byte(output))
									}
								}

								// Clean up
								os.Remove(cmdFile)
							}
						}
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
				remaining := lineBuffer.String()
				if len(remaining) > 0 {
					cleanLine := i.stripANSI(remaining)
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
