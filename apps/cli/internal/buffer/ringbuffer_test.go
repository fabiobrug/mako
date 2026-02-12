package buffer

import (
	"testing"
)

func TestRingBuffer(t *testing.T) {
	rb := NewRingBuffer(3)

	// Test 1: Add lines without overflow
	rb.Write("line 1")
	rb.Write("line 2")

	lines := rb.GetAll()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "line 1" || lines[1] != "line 2" {
		t.Errorf("Lines don't match: %v", lines)
	}

	// Test 2: Overflow (circular behavior)
	rb.Write("line 3")
	rb.Write("line 4")

	lines = rb.GetAll()
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "line 2" || lines[1] != "line 3" || lines[2] != "line 4" {
		t.Errorf("Expected [line 2, line 3, line 4], got %v", lines)
	}

	// Test 3: GetLines with limit
	recent := rb.GetLines(2)
	if len(recent) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(recent))
	}
	if recent[0] != "line 3" || recent[1] != "line 4" {
		t.Errorf("Expected last 2 lines [line 3, line 4], got %v", recent)
	}

	// Test 4: Clear
	rb.Clear()
	if rb.Count() != 0 {
		t.Errorf("Expected 0 lines after clear, got %d", rb.Count())
	}

}
