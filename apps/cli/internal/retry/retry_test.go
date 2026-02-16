package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.MaxAttempts != 3 {
		t.Errorf("Expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay != 100*time.Millisecond {
		t.Errorf("Expected InitialDelay=100ms, got %v", cfg.InitialDelay)
	}
	if cfg.MaxDelay != 5*time.Second {
		t.Errorf("Expected MaxDelay=5s, got %v", cfg.MaxDelay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("Expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}

func TestDoSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()
	attempts := 0

	err := Do(ctx, cfg, func() error {
		attempts++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestDoRetryAndSuccess(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.InitialDelay = 10 * time.Millisecond
	attempts := 0

	err := Do(ctx, cfg, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestDoMaxAttemptsExceeded(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.InitialDelay = 10 * time.Millisecond
	attempts := 0

	err := Do(ctx, cfg, func() error {
		attempts++
		return errors.New("persistent error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != cfg.MaxAttempts {
		t.Errorf("Expected %d attempts, got %d", cfg.MaxAttempts, attempts)
	}
}

func TestDoContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := DefaultConfig()
	cfg.InitialDelay = 100 * time.Millisecond
	attempts := 0

	// Cancel context after first attempt
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := Do(ctx, cfg, func() error {
		attempts++
		return errors.New("error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts > 2 {
		t.Errorf("Expected at most 2 attempts, got %d", attempts)
	}
}

func TestDoNonRetryableError(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.InitialDelay = 10 * time.Millisecond
	cfg.RetryableErrors = func(err error) bool {
		return err.Error() != "non-retryable"
	}
	attempts := 0

	err := Do(ctx, cfg, func() error {
		attempts++
		return errors.New("non-retryable")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestDoWithResult(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.InitialDelay = 10 * time.Millisecond
	attempts := 0

	result, err := DoWithResult(ctx, cfg, func() (string, error) {
		attempts++
		if attempts < 2 {
			return "", errors.New("temporary error")
		}
		return "success", nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got '%s'", result)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestDoWithResultError(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.MaxAttempts = 2
	cfg.InitialDelay = 10 * time.Millisecond

	result, err := DoWithResult(ctx, cfg, func() (int, error) {
		return 0, errors.New("persistent error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != 0 {
		t.Errorf("Expected zero value result, got %d", result)
	}
}

func TestCalculateDelay(t *testing.T) {
	cfg := &Config{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
		Jitter:       false,
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{1, 100 * time.Millisecond},
		{2, 200 * time.Millisecond},
		{3, 400 * time.Millisecond},
		{4, 800 * time.Millisecond},
		{5, 1 * time.Second}, // Capped at MaxDelay
		{6, 1 * time.Second}, // Capped at MaxDelay
	}

	for _, tt := range tests {
		got := cfg.calculateDelay(tt.attempt)
		if got != tt.want {
			t.Errorf("calculateDelay(%d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}

func TestCalculateDelayWithJitter(t *testing.T) {
	cfg := &Config{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
	}

	// Run multiple times to test jitter varies
	delays := make([]time.Duration, 5)
	for i := 0; i < 5; i++ {
		delays[i] = cfg.calculateDelay(2)
	}

	// Check that all delays are within reasonable range
	for _, delay := range delays {
		// 200ms * (1 ± 0.2) = 160ms to 240ms
		if delay < 160*time.Millisecond || delay > 240*time.Millisecond {
			t.Errorf("Delay %v out of expected range [160ms, 240ms]", delay)
		}
	}
}

func TestAggressiveConfig(t *testing.T) {
	cfg := AggressiveConfig()

	if cfg.MaxAttempts != 5 {
		t.Errorf("Expected MaxAttempts=5, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay != 50*time.Millisecond {
		t.Errorf("Expected InitialDelay=50ms, got %v", cfg.InitialDelay)
	}
}

func BenchmarkDo(b *testing.B) {
	ctx := context.Background()
	cfg := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Do(ctx, cfg, func() error {
			return nil
		})
	}
}

func BenchmarkDoWithRetries(b *testing.B) {
	ctx := context.Background()
	cfg := DefaultConfig()
	cfg.InitialDelay = 1 * time.Microsecond // Very short for benchmark

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		attempts := 0
		_ = Do(ctx, cfg, func() error {
			attempts++
			if attempts < 2 {
				return errors.New("retry")
			}
			return nil
		})
	}
}
