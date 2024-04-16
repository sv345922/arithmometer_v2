package entities

import (
	"time"
)

// Тип дл работы с базой данных
type Task struct {
	NodeId   uint64    `json:"node_id"`  // id, образующего задачу нода
	X        float64   `json:"x"`        // операнд X
	XReady   bool      `json:"x_ready"`  // готовность операнда X
	Y        float64   `json:"y"`        // операнд Y
	YReady   bool      `json:"y_ready"`  // готовность операнда Y
	CalcId   uint64    `json:"calc_id"`  // id вычислителя задачи
	Deadline time.Time `json:"deadline"` // дедлайн задачи
	Duration int       `json:"duration"` // длительность операции
}
