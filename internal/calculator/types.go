package calculator

import (
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

type Calculator struct {
	Id   int
	Task *entities.MessageTask
	Ch   chan entities.MessageResult
}

var URL = "http://127.0.0.1:" + configs.Port
