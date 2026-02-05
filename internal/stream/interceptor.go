package stream

import (
	"bytes"
	"io"
	"regexp"
	"strings"

	"github.com/fabiobrug/mako.git/internal/buffer"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

type Interceptor struct {
	buffer *buffer.RingBuffer
}

func NewInterceptor(bufferSize int) *Interceptor {
	return &Interceptor{
		buffer: buffer.NewRingBuffer(bufferSize),
	}
}

func (i *Interceptor) Tee(dst io.Writer, src io.Reader) error {
	buf := make([]byte, 1024)
	lineBuffer := bytes.NewBuffer(nil)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, writeErr := dst.Write(buf[:n]); writeErr != nil {
				return writeErr
			}

			lineBuffer.Write(buf[:n])

			data := lineBuffer.Bytes()
			lastNewline := bytes.LastIndexByte(data, '\n')

			if lastNewline >= 0 {
				lines := bytes.Split(data[:lastNewline+1], []byte{'\n'})
				for _, line := range lines {
					if len(line) == 0 {
						continue
					}
					cleanLine := i.stripANSI(string(line))
					cleanLine = strings.TrimSpace(cleanLine)
					if len(cleanLine) > 0 {
						i.buffer.Write(cleanLine)
					}
				}

				lineBuffer.Reset()
				if lastNewline+1 < len(data) {
					lineBuffer.Write(data[lastNewline+1:])
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
