package getResult

import (
	"encoding/json"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

// Возвращает json массив байт выражений, у которых соответсвует userId
func GetExpressions(userId uint64, expessions map[uint64]*entities.Expression) ([]byte, error) {
	// Получаем id выражений (ключи)
	exprIds := make([]uint64, 0, len(expessions))
	for id := range expessions {
		exprIds = append(exprIds, id)
	}
	resultExpressions := make(map[uint64]*entities.Expression)
	for _, id := range exprIds {
		if expression, ok := expessions[id]; ok && expression.UserId == userId {
			resultExpressions[id] = expression
		}
	}

	result, err := json.Marshal(resultExpressions)
	if err != nil {
		return nil, err
	}
	return result, nil
}
