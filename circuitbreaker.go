package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// State — состояния Circuit Breaker
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// Config — конфигурация Circuit Breaker
type Config struct {
	Name        string
	MaxRequests uint32           // Максимальное количество ошибок
	Interval    time.Duration    // Время ожидания перед переходом в HalfOpen
	Timeout     time.Duration    // Таймаут на выполнение запроса
	IsError     func(error) bool // Функция, определяющая, является ли ошибка критичной
}

// CircuitBreaker — реализация паттерна Circuit Breaker
type CircuitBreaker struct {
	cfg      Config
	state    State
	mu       sync.Mutex
	failures uint32
	timer    *time.Timer
}

// NewCircuitBreaker — создаёт новый Circuit Breaker
func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	if cfg.IsError == nil {
		cfg.IsError = func(err error) bool {
			return err != nil
		}
	}

	return &CircuitBreaker{
		cfg:   cfg,
		state: Closed,
	}
}

// Execute — выполняет функцию с Circuit Breaker
func (cb *CircuitBreaker) Execute(fn func() (string, error)) (string, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == Open {
		return "", fmt.Errorf("circuit breaker is open: %w", ErrCircuitBreakerOpen)
	}

	if cb.state == HalfOpen {
		// В HalfOpen режиме разрешён только один тестовый запрос
		cb.failures = 0
	}

	result, err := fn()
	if err != nil && cb.cfg.IsError(err) {
		cb.failures++
		if cb.failures >= cb.cfg.MaxRequests {
			cb.state = Open
			cb.mu.Unlock()
			cb.timer = time.AfterFunc(cb.cfg.Interval, func() {
				cb.mu.Lock()
				cb.state = HalfOpen
				cb.failures = 0
				cb.mu.Unlock()
			})
			cb.mu.Lock()
		}
		return "", err
	}

	// Сброс ошибок при успешном запросе
	cb.failures = 0

	// Если мы в HalfOpen и запрос успешен, возвращаемся в Closed
	if cb.state == HalfOpen {
		cb.state = Closed
		if cb.timer != nil {
			cb.timer.Stop()
		}
	}

	return result, nil
}
