package entities

import "sync"

// Expression Выражение
type Expression struct {
	Id         uint64  `json:"id"` // Id запроса клиента
	UserId     uint64  `json:"user_id"`
	UserTask   string  `json:"userTask"` // Задание клиента
	ResultExpr float64 `json:"result"`   // Результат,
	Status     string  `json:"status"`   // ""/"done"/"zero division"
	RootId     uint64  `json:"rootId"`   // Id корневого узла
	mu         sync.RWMutex
}

func (e *Expression) SetStatus(status string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Status = status
}
func (e *Expression) SetRoot(id uint64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.RootId = id
}
