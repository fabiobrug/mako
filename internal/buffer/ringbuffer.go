package buffer

import (
	"strings"
	"sync"
)

type RingBuffer struct {
	data  []string
	size  int
	head  int
	count int
	mu    sync.RWMutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:  make([]string, size),
		size:  size,
		head:  0,
		count: 0,
	}
}

func (rb *RingBuffer) Write(line string) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	line = strings.TrimRight(line, "\r\n")

	rb.data[rb.head] = line

	rb.head = (rb.head + 1) % rb.size

	if rb.count < rb.size {
		rb.count++
	}
}

func (rb *RingBuffer) GetLines(n int) []string {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if n > rb.count {
		n = rb.count
	}

	if n == 0 {
		return []string{}
	}

	result := make([]string, n)

	var start int
	if rb.count < rb.size {
		start = 0
	} else {
		start = rb.head
	}

	for i := 0; i < n; i++ {
		idx := (start + rb.count - n + i) % rb.size
		result[i] = rb.data[idx]
	}

	return result
}

func (rb *RingBuffer) GetAll() []string {
	return rb.GetLines(rb.count)
}

func (rb *RingBuffer) Clear() {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.head = 0
	rb.count = 0
}

func (rb *RingBuffer) Count() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	return rb.count
}
