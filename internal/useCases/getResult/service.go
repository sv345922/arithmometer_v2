package getResult

import (
	"encoding/json"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

// Возвращает json массив байт выражений, у которых соответсвует userId
func GetExpressions(userId uint64, expressions map[uint64]*entities.Expression) ([]byte, error) {
	// Получаем id выражений (ключи)
	exprIds := make([]uint64, 0, len(expressions))
	for id := range expressions {
		exprIds = append(exprIds, id)
	}
	resultExpressions := make(map[uint64]*entities.Expression)
	for _, id := range exprIds {
		// Проверка на существование выражения и соответсвие пользователю
		if expression, ok := expressions[id]; ok && expression.UserId == userId {
			resultExpressions[id] = expression
		}
	}

	result, err := json.Marshal(resultExpressions)
	if err != nil {
		return nil, err
	}
	return result, nil
}
