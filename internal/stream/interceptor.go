package stream

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/fabiobrug/mako.git/internal/buffer"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-z]`)

type Interceptor struct {
	buffer *buffer.RingBuffer
}

func (i *Interceptor) Tee(dst io.Writer, src io.Reader) error {
	scanner := bufio.NewScanner(src)

	for scanner.Scan() {
		line := scanner.Text()

		if _, err := dst.Write([]byte(line + "\n")); err != nil {
			return err
		}

		cleanLine := i.stripANSI(line)
		if len(strings.TrimSpace(cleanLine)) > 0 {
			i.buffer.Write(cleanLine)
		}
	}

	return scanner.Err()
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
