package entities

import (
	"sync"
)

// Node - узел выражения
type Node struct {
	Id           uint64  `json:"nodeId"`
	ExpressionId uint64  `json:"expressionId"` // id выражения
	Op           string  `json:"op"`           // оператор
	X            uint64  `json:"x"`
	Y            uint64  `json:"y"`          // потомки
	Val          float64 `json:"val"`        // значение узла
	Sheet        bool    `json:"sheet"`      // флаг листа
	Calculated   bool    `json:"calculated"` // флаг вычисленного узла
	Parent       uint64  `json:"parent"`     // узел родитель
	mu           sync.RWMutex
}

func (n *Node) IsSheet() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Sheet
}
