package calculator

import (
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

type Calculator struct {
	Id   uint64
	Task *entities.MessageTask
}

var URL = "http://127.0.0.1:" + configs.Port
