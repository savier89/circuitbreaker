# Circuit Breaker

## 🚧 Реализация паттерна Circuit Breaker на Go

`circuitbreaker` — это простая и эффективная реализация паттерна [Circuit Breaker](https://learn.microsoft.com/en-us/azure/architecture/reference-architectures/serverless/event-hub-stream-analytics-cosmos-db#circuit-breaker-pattern) на языке Go, предназначенная для повышения отказоустойчивости ваших сервисов. Библиотека помогает предотвращать каскадные сбои, обрабатывать временные ошибки и улучшать стабильность распределённых систем.

## 📌 Особенности

- Поддержка синхронных и асинхронных вызовов
- Настройка порога ошибок и времени восстановления
- Поддержка пользовательских fallback-функций
- Лёгкая и понятная в использовании
- Совместимость с Go 1.18+

## 📦 Установка

Установите библиотеку с помощью команды `go get`:

```bash
go get github.com/savier89/circuitbreaker
```

Импортируйте в свой проект:

```go
import "github.com/savier89/circuitbreaker"
```

## 🧪 Пример использования

### 1. Базовое использование

```go
package main

import (
    "errors"
    "fmt"
    "github.com/savier89/circuitbreaker"
)

func unreliableFunc() (string, error) {
    // Имитация нестабильного вызова
    return "", errors.New("service unavailable")
}

func main() {
    cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.Settings{
        Name:             "example-cb",
        FailureThreshold: 3,
        RecoveryTime:     10, // в секундах
    })

    for i := 0; i < 6; i++ {
        result, err := cb.Execute(unreliableFunc)
        if err != nil {
            fmt.Printf("Attempt %d: Circuit breaker error: %s\n", i+1, err)
            continue
        }
        fmt.Println("Result:", result)
    }
}
```

### 2. Добавление fallback функции

```go
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.Settings{
    Name:             "fallback-cb",
    FailureThreshold: 2,
    RecoveryTime:     30,
    Fallback: func() (interface{}, error) {
        return "Fallback response", nil
    },
})

result, err := cb.Execute(func() (interface{}, error) {
    return "", errors.New("external service failed")
})

fmt.Println("Result:", result) // Выведет: Fallback response
```

### 3. Настройка логирования и событий

Вы можете подписаться на события цепи (например, переход в состояние "открыто" или "закрыто") для логирования или мониторинга.

```go
cb.OnStateChange = func(name string, from, to circuitbreaker.State) {
    fmt.Printf("Circuit [%s] state changed: %s → %s\n", name, from, to)
}
```

## 🧰 Поддерживаемые версии

- Go 1.18 и выше
- Поддержка generics (начиная с Go 1.18)

## 📚 Дополнительная информация

Паттерн Circuit Breaker особенно полезен в микросервисных архитектурах при работе с внешними API, базами данных и другими сторонними сервисами, где возможны временные сбои. Библиотека позволяет гибко настраивать поведение и встраивать отказоустойчивость в ваши приложения.

## 🤝 Участие в проекте

Любые предложения, баг-репорты и пул-реквесты приветствуются! Пожалуйста, используйте [GitHub Issues](https://github.com/savier89/circuitbreaker/issues) для обратной связи.

## 📄 Лицензия

Этот проект лицензируется под [MIT License](LICENSE).

---

Если у вас есть специфические особенности реализации (например, поддержка контекста, таймеров, middleware и т.д.), я могу дополнить README под них.