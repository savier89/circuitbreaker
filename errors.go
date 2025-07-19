package circuitbreaker

import "errors"

// Ошибки
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
	ErrExecutionTimeout   = errors.New("execution timeout")
)
