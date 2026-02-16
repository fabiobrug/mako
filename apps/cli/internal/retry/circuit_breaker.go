package retry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state
type State int

const (
	// StateClosed means circuit is closed and requests are allowed
	StateClosed State = iota
	// StateOpen means circuit is open and requests are blocked
	StateOpen
	// StateHalfOpen means circuit is testing if backend has recovered
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

var (
	// ErrCircuitOpen is returned when the circuit breaker is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTooManyRequests is returned when too many requests are in half-open state
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	// MaxFailures is the number of failures before opening the circuit
	MaxFailures uint32
	// Timeout is how long to wait in open state before trying half-open
	Timeout time.Duration
	// MaxRequests is the maximum number of requests allowed in half-open state
	MaxRequests uint32
	// OnStateChange is called when the state changes
	OnStateChange func(from, to State)
}

// DefaultCircuitBreakerConfig returns sensible defaults
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		MaxFailures: 5,
		Timeout:     30 * time.Second,
		MaxRequests: 1,
		OnStateChange: func(from, to State) {
			// Default: no-op
		},
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config *CircuitBreakerConfig
	mu     sync.RWMutex

	state            State
	failures         uint32
	successes        uint32
	lastStateChange  time.Time
	lastFailureTime  time.Time
	halfOpenRequests uint32
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}

	return &CircuitBreaker{
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Execute runs the given operation through the circuit breaker
func (cb *CircuitBreaker) Execute(ctx context.Context, operation func() error) error {
	// Check if we can execute
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	// Execute the operation
	err := operation()

	// Record the result
	cb.afterRequest(err)

	return err
}

// beforeRequest checks if the request should be allowed
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		// Check if we should transition to half-open
		if time.Since(cb.lastStateChange) > cb.config.Timeout {
			cb.setState(StateHalfOpen)
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// Limit concurrent requests in half-open state
		if cb.halfOpenRequests >= cb.config.MaxRequests {
			return ErrTooManyRequests
		}
		cb.halfOpenRequests++
		return nil

	default:
		return ErrCircuitOpen
	}
}

// afterRequest records the result of a request
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

// onSuccess handles a successful request
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		// Reset failure count on success
		cb.failures = 0

	case StateHalfOpen:
		cb.successes++
		cb.halfOpenRequests--

		// If we've had enough successes, close the circuit
		if cb.successes >= cb.config.MaxRequests {
			cb.setState(StateClosed)
			cb.failures = 0
			cb.successes = 0
		}
	}
}

// onFailure handles a failed request
func (cb *CircuitBreaker) onFailure() {
	cb.lastFailureTime = time.Now()
	cb.failures++

	switch cb.state {
	case StateClosed:
		// Open circuit if we've exceeded max failures
		if cb.failures >= cb.config.MaxFailures {
			cb.setState(StateOpen)
		}

	case StateHalfOpen:
		// Go back to open on any failure in half-open state
		cb.halfOpenRequests--
		cb.setState(StateOpen)
		cb.successes = 0
	}
}

// setState changes the circuit breaker state
func (cb *CircuitBreaker) setState(newState State) {
	oldState := cb.state
	if oldState == newState {
		return
	}

	cb.state = newState
	cb.lastStateChange = time.Now()

	if cb.config.OnStateChange != nil {
		// Call state change handler without holding lock
		go cb.config.OnStateChange(oldState, newState)
	}
}

// State returns the current circuit breaker state
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Stats returns circuit breaker statistics
type CircuitBreakerStats struct {
	State            State
	Failures         uint32
	Successes        uint32
	LastStateChange  time.Time
	LastFailureTime  time.Time
	HalfOpenRequests uint32
}

func (cb *CircuitBreaker) Stats() CircuitBreakerStats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return CircuitBreakerStats{
		State:            cb.state,
		Failures:         cb.failures,
		Successes:        cb.successes,
		LastStateChange:  cb.lastStateChange,
		LastFailureTime:  cb.lastFailureTime,
		HalfOpenRequests: cb.halfOpenRequests,
	}
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateClosed)
	cb.failures = 0
	cb.successes = 0
	cb.halfOpenRequests = 0
}

// ResilientExecutor combines retry logic with circuit breaker
type ResilientExecutor struct {
	Retry          *Config
	CircuitBreaker *CircuitBreaker
	// Keep lowercase aliases for backward compatibility within package
	retry          *Config
	circuitBreaker *CircuitBreaker
}

// NewResilientExecutor creates a new resilient executor
func NewResilientExecutor(retryConfig *Config, cb *CircuitBreaker) *ResilientExecutor {
	if retryConfig == nil {
		retryConfig = DefaultConfig()
	}
	if cb == nil {
		cb = NewCircuitBreaker(DefaultCircuitBreakerConfig())
	}

	return &ResilientExecutor{
		Retry:          retryConfig,
		CircuitBreaker: cb,
		retry:          retryConfig,
		circuitBreaker: cb,
	}
}

// Execute runs an operation with both retry and circuit breaker
func (re *ResilientExecutor) Execute(ctx context.Context, operation func() error) error {
	return Do(ctx, re.retry, func() error {
		return re.circuitBreaker.Execute(ctx, operation)
	})
}

// CircuitBreakerStats returns the circuit breaker stats
func (re *ResilientExecutor) CircuitBreakerStats() CircuitBreakerStats {
	return re.CircuitBreaker.Stats()
}

// ResetCircuitBreaker resets the circuit breaker
func (re *ResilientExecutor) ResetCircuitBreaker() {
	re.CircuitBreaker.Reset()
}

// IsRetryableHTTPError determines if an HTTP error is retryable
func IsRetryableHTTPError(statusCode int) bool {
	// Retry on server errors and rate limiting
	switch statusCode {
	case 408, // Request Timeout
		429, // Too Many Requests
		500, // Internal Server Error
		502, // Bad Gateway
		503, // Service Unavailable
		504: // Gateway Timeout
		return true
	default:
		return false
	}
}

// RetryableErrorFunc creates a retry config that checks HTTP status codes
func RetryableErrorFunc(isRetryable func(error) bool) func(error) bool {
	return func(err error) bool {
		if err == nil {
			return false
		}

		// Don't retry circuit breaker errors
		if errors.Is(err, ErrCircuitOpen) || errors.Is(err, ErrTooManyRequests) {
			return false
		}

		// Use custom function if provided
		if isRetryable != nil {
			return isRetryable(err)
		}

		// By default, retry most errors except explicit non-retryable ones
		return true
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
