package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(nil)

	if cb == nil {
		t.Fatal("Expected non-nil circuit breaker")
	}
	if cb.State() != StateClosed {
		t.Errorf("Expected initial state Closed, got %v", cb.State())
	}
}

func TestCircuitBreakerClosedState(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	ctx := context.Background()

	err := cb.Execute(ctx, func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed, got %v", cb.State())
	}
}

func TestCircuitBreakerOpensAfterFailures(t *testing.T) {
	config := &CircuitBreakerConfig{
		MaxFailures: 3,
		Timeout:     1 * time.Second,
		MaxRequests: 1,
	}
	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	// Trigger failures
	for i := 0; i < 3; i++ {
		_ = cb.Execute(ctx, func() error {
			return errors.New("failure")
		})
	}

	if cb.State() != StateOpen {
		t.Errorf("Expected state Open after failures, got %v", cb.State())
	}

	// Next request should be blocked
	err := cb.Execute(ctx, func() error {
		return nil
	})

	if !errors.Is(err, ErrCircuitOpen) {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	config := &CircuitBreakerConfig{
		MaxFailures: 2,
		Timeout:     100 * time.Millisecond,
		MaxRequests: 1,
	}
	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	// Trigger failures to open circuit
	for i := 0; i < 2; i++ {
		_ = cb.Execute(ctx, func() error {
			return errors.New("failure")
		})
	}

	if cb.State() != StateOpen {
		t.Fatalf("Expected state Open, got %v", cb.State())
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Next request should transition to half-open
	err := cb.Execute(ctx, func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error in half-open, got %v", err)
	}

	// Circuit should now be closed
	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed after success in half-open, got %v", cb.State())
	}
}

func TestCircuitBreakerHalfOpenFails(t *testing.T) {
	config := &CircuitBreakerConfig{
		MaxFailures: 2,
		Timeout:     100 * time.Millisecond,
		MaxRequests: 1,
	}
	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	// Open circuit
	for i := 0; i < 2; i++ {
		_ = cb.Execute(ctx, func() error {
			return errors.New("failure")
		})
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Fail in half-open
	_ = cb.Execute(ctx, func() error {
		return errors.New("failure")
	})

	// Circuit should reopen
	if cb.State() != StateOpen {
		t.Errorf("Expected state Open after failure in half-open, got %v", cb.State())
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	config := &CircuitBreakerConfig{
		MaxFailures: 2,
		Timeout:     1 * time.Second,
		MaxRequests: 1,
	}
	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	// Open circuit
	for i := 0; i < 2; i++ {
		_ = cb.Execute(ctx, func() error {
			return errors.New("failure")
		})
	}

	if cb.State() != StateOpen {
		t.Fatalf("Expected state Open, got %v", cb.State())
	}

	// Reset
	cb.Reset()

	if cb.State() != StateClosed {
		t.Errorf("Expected state Closed after reset, got %v", cb.State())
	}

	// Should work now
	err := cb.Execute(ctx, func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error after reset, got %v", err)
	}
}

func TestCircuitBreakerStats(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	ctx := context.Background()

	// Execute some operations
	cb.Execute(ctx, func() error { return nil })
	cb.Execute(ctx, func() error { return errors.New("fail") })

	stats := cb.Stats()

	if stats.State != StateClosed {
		t.Errorf("Expected state Closed, got %v", stats.State)
	}
	if stats.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", stats.Failures)
	}
}

func TestCircuitBreakerStateChange(t *testing.T) {
	stateChanges := []struct{ from, to State }{}
	config := &CircuitBreakerConfig{
		MaxFailures: 2,
		Timeout:     100 * time.Millisecond,
		MaxRequests: 1,
		OnStateChange: func(from, to State) {
			stateChanges = append(stateChanges, struct{ from, to State }{from, to})
		},
	}
	cb := NewCircuitBreaker(config)
	ctx := context.Background()

	// Trigger state change to Open
	for i := 0; i < 2; i++ {
		_ = cb.Execute(ctx, func() error {
			return errors.New("failure")
		})
	}

	// Wait for callback
	time.Sleep(50 * time.Millisecond)

	if len(stateChanges) != 1 {
		t.Fatalf("Expected 1 state change, got %d", len(stateChanges))
	}
	if stateChanges[0].from != StateClosed || stateChanges[0].to != StateOpen {
		t.Errorf("Expected Closed->Open, got %v->%v", stateChanges[0].from, stateChanges[0].to)
	}
}

func TestResilientExecutor(t *testing.T) {
	retryConfig := &Config{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}
	cbConfig := &CircuitBreakerConfig{
		MaxFailures: 5,
		Timeout:     1 * time.Second,
		MaxRequests: 1,
	}

	executor := NewResilientExecutor(retryConfig, NewCircuitBreaker(cbConfig))
	ctx := context.Background()
	attempts := 0

	err := executor.Execute(ctx, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestResilientExecutorCombination(t *testing.T) {
	retryConfig := &Config{
		MaxAttempts:  2,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}
	cbConfig := &CircuitBreakerConfig{
		MaxFailures: 3,
		Timeout:     1 * time.Second,
		MaxRequests: 1,
	}

	executor := NewResilientExecutor(retryConfig, NewCircuitBreaker(cbConfig))
	ctx := context.Background()

	// Test that retry works
	attempts := 0
	err := executor.Execute(ctx, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestIsRetryableHTTPError(t *testing.T) {
	tests := []struct {
		statusCode int
		want       bool
	}{
		{200, false},
		{400, false},
		{401, false},
		{403, false},
		{408, true},  // Request Timeout
		{429, true},  // Too Many Requests
		{500, true},  // Internal Server Error
		{502, true},  // Bad Gateway
		{503, true},  // Service Unavailable
		{504, true},  // Gateway Timeout
	}

	for _, tt := range tests {
		got := IsRetryableHTTPError(tt.statusCode)
		if got != tt.want {
			t.Errorf("IsRetryableHTTPError(%d) = %v, want %v", tt.statusCode, got, tt.want)
		}
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(999), "unknown"},
	}

	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("State(%v).String() = %s, want %s", tt.state, got, tt.want)
		}
	}
}

func BenchmarkCircuitBreakerExecute(b *testing.B) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cb.Execute(ctx, func() error {
			return nil
		})
	}
}
