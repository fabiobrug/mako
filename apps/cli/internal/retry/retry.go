package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Config holds retry configuration
type Config struct {
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int
	// InitialDelay is the initial delay between retries
	InitialDelay time.Duration
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	// Multiplier is the backoff multiplier
	Multiplier float64
	// Jitter adds randomness to delays to prevent thundering herd
	Jitter bool
	// RetryableErrors is a function that determines if an error is retryable
	RetryableErrors func(error) bool
}

// DefaultConfig returns sensible defaults for API retries
func DefaultConfig() *Config {
	return &Config{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableErrors: func(err error) bool {
			// By default, retry all errors
			// Specific implementations can override this
			return true
		},
	}
}

// AggressiveConfig returns a more aggressive retry configuration
func AggressiveConfig() *Config {
	return &Config{
		MaxAttempts:  5,
		InitialDelay: 50 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableErrors: func(err error) bool {
			return true
		},
	}
}

// Do executes a function with retries using exponential backoff
func Do(ctx context.Context, cfg *Config, operation func() error) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		// Attempt the operation
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if cfg.RetryableErrors != nil && !cfg.RetryableErrors(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		// Don't sleep on last attempt
		if attempt >= cfg.MaxAttempts {
			break
		}

		// Calculate delay with exponential backoff
		delay := cfg.calculateDelay(attempt)

		// Wait before retry
		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return fmt.Errorf("operation cancelled during retry: %w", ctx.Err())
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", cfg.MaxAttempts, lastErr)
}

// DoWithResult executes a function that returns a result with retries
func DoWithResult[T any](ctx context.Context, cfg *Config, operation func() (T, error)) (T, error) {
	var result T
	var lastErr error

	if cfg == nil {
		cfg = DefaultConfig()
	}

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return result, fmt.Errorf("operation cancelled: %w", ctx.Err())
		default:
		}

		// Attempt the operation
		res, err := operation()
		if err == nil {
			return res, nil
		}

		lastErr = err

		// Check if error is retryable
		if cfg.RetryableErrors != nil && !cfg.RetryableErrors(err) {
			return result, fmt.Errorf("non-retryable error: %w", err)
		}

		// Don't sleep on last attempt
		if attempt >= cfg.MaxAttempts {
			break
		}

		// Calculate delay with exponential backoff
		delay := cfg.calculateDelay(attempt)

		// Wait before retry
		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return result, fmt.Errorf("operation cancelled during retry: %w", ctx.Err())
		}
	}

	return result, fmt.Errorf("operation failed after %d attempts: %w", cfg.MaxAttempts, lastErr)
}

// calculateDelay computes the delay for a given attempt using exponential backoff
func (c *Config) calculateDelay(attempt int) time.Duration {
	// Calculate exponential backoff: initialDelay * (multiplier ^ (attempt - 1))
	delay := float64(c.InitialDelay) * math.Pow(c.Multiplier, float64(attempt-1))

	// Apply maximum delay cap
	if delay > float64(c.MaxDelay) {
		delay = float64(c.MaxDelay)
	}

	// Add jitter if enabled (±20% randomness)
	if c.Jitter {
		jitter := delay * 0.2 * (2*rand() - 1) // Random value between -20% and +20%
		delay += jitter
	}

	return time.Duration(delay)
}

// Simple pseudo-random function for jitter
func rand() float64 {
	// Use current time nanoseconds for simple randomness
	// This is sufficient for jitter and avoids importing crypto/rand
	ns := time.Now().UnixNano()
	return float64(ns%1000) / 1000.0
}
